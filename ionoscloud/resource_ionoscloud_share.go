package ionoscloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceIonosCloudShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceIonosCloudShareCreate,
		Read:   resourceIonosCloudShareRead,
		Update: resourceIonosCloudShareUpdate,
		Delete: resourceIonosCloudShareDelete,
		Schema: map[string]*schema.Schema{
			"edit_privilege": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"share_privilege": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceIonosCloudShareCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	request := ionoscloud.Share{
		Properties: ionoscloud.ShareProperties{},
	}

	tempSharePrivilege := d.Get("edit_privilege").(bool)
	request.Properties.SharePrivilege = &tempSharePrivilege
	tempEditPrivilege := d.Get("share_privilege").(bool)
	request.Properties.EditPrivilege = &tempEditPrivilege

	share, err := client.AddShare(d.Get("group_id").(string), d.Get("resource_id").(string), request)

	log.Printf("[DEBUG] SHARE ID: %s", share.ID)

	if err != nil {
		return fmt.Errorf("An error occured while creating a share: %s", err)
	}
	d.SetId(share.ID)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, share.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceIonosCloudShareRead(d, meta)
}

func resourceIonosCloudShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	share, err := client.GetShare(d.Get("group_id").(string), d.Get("resource_id").(string))

	if err != nil {
		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a Share ID %s %s", d.Id(), err)
	}

	d.Set("edit_privilege", share.Properties.EditPrivilege)
	d.Set("share_privilege", share.Properties.SharePrivilege)
	return nil
}

func resourceIonosCloudShareUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	tempSharePrivilege := d.Get("share_privilege").(bool)
	tempEditPrivilege := d.Get("edit_privilege").(bool)
	shareReq := ionoscloud.Share{
		Properties: ionoscloud.ShareProperties{
			EditPrivilege:  &tempEditPrivilege,
			SharePrivilege: &tempSharePrivilege,
		},
	}

	share, err := client.UpdateShare(d.Get("group_id").(string), d.Get("resource_id").(string), shareReq)
	if err != nil {
		return fmt.Errorf("An error occured while patching a share ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, share.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceIonosCloudShareRead(d, meta)
}

func resourceIonosCloudShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	resp, err := client.DeleteShare(d.Id(), d.Get("resource_id").(string))
	if err != nil {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		resp, err = client.DeleteShare(d.Id(), d.Get("resource_id").(string))
		if err != nil {
			if apiError, ok := err.(ionoscloud.ApiError); ok {
				if apiError.HttpStatusCode() != 404 {
					return fmt.Errorf("An error occured while deleting a share %s %s", d.Id(), err)
				}
			}
		}
	}

	// Wait, catching any errors
	if resp.Get("Location") != "" {
		_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
		if errState != nil {
			return errState
		}
	}

	d.SetId("")
	return nil
}
