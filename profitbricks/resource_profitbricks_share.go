package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
	"log"
	"time"
)

func resourceProfitBricksShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksShareCreate,
		Read:   resourceProfitBricksShareRead,
		Update: resourceProfitBricksShareUpdate,
		Delete: resourceProfitBricksShareDelete,
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
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProfitBricksShareCreate(d *schema.ResourceData, meta interface{}) error {
	request := profitbricks.Share{
		Properties: profitbricks.ShareProperties{},
	}

	tempSharePrivilege := d.Get("edit_privilege").(bool)
	request.Properties.SharePrivilege = &tempSharePrivilege
	tempEditPrivilege := d.Get("share_privilege").(bool)
	request.Properties.EditPrivilege = &tempEditPrivilege

	share := profitbricks.AddShare(request, d.Get("group_id").(string), d.Get("resource_id").(string))

	log.Printf("[DEBUG] SHARE ID: %s", share.Id)

	if share.StatusCode > 299 {
		return fmt.Errorf("An error occured while creating a share: %s", share.Response)
	}

	err := waitTillProvisioned(meta, share.Headers.Get("Location"))
	if err != nil {
		return err
	}
	d.SetId(share.Id)
	return resourceProfitBricksShareRead(d, meta)
}

func resourceProfitBricksShareRead(d *schema.ResourceData, meta interface{}) error {
	share := profitbricks.GetShare(d.Get("group_id").(string), d.Get("resource_id").(string))

	if share.StatusCode > 299 {
		if share.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a Share ID %s %s", d.Id(), share.Response)
	}

	d.Set("edit_privilege", share.Properties.EditPrivilege)
	d.Set("share_privilege", share.Properties.SharePrivilege)
	return nil
}

func resourceProfitBricksShareUpdate(d *schema.ResourceData, meta interface{}) error {
	tempSharePrivilege := d.Get("share_privilege").(bool)
	tempEditPrivilege := d.Get("edit_privilege").(bool)
	shareReq := profitbricks.Share{
		Properties: profitbricks.ShareProperties{
			EditPrivilege:  &tempEditPrivilege,
			SharePrivilege: &tempSharePrivilege,
		},
	}

	share := profitbricks.UpdateShare(d.Get("group_id").(string), d.Get("resource_id").(string), shareReq)
	if share.StatusCode > 299 {
		return fmt.Errorf("An error occured while patching a share ID %s %s", d.Id(), share.Response)
	}
	err := waitTillProvisioned(meta, share.Headers.Get("Location"))
	if err != nil {
		return err
	}
	return resourceProfitBricksShareRead(d, meta)
}

func resourceProfitBricksShareDelete(d *schema.ResourceData, meta interface{}) error {
	resp := profitbricks.DeleteShare(d.Id(), d.Get("resource_id").(string))
	if resp.StatusCode > 299 {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		resp = profitbricks.DeleteShare(d.Id(), d.Get("resource_id").(string))
		if resp.StatusCode > 299 && resp.StatusCode != 404 {
			return fmt.Errorf("An error occured while deleting a share %s %s", d.Id(), string(resp.Body))
		}
	}

	if resp.Headers.Get("Location") != "" {
		err := waitTillProvisioned(meta, resp.Headers.Get("Location"))
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}
