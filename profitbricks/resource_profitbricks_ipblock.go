package profitbricks

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceProfitBricksIPBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksIPBlockCreate,
		Read:   resourceProfitBricksIPBlockRead,
		Update: resourceProfitBricksIPBlockUpdate,
		Delete: resourceProfitBricksIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksIPBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	ipblock := &profitbricks.IPBlock{
		Properties: profitbricks.IPBlockProperties{
			Size:     d.Get("size").(int),
			Location: d.Get("location").(string),
			Name:     d.Get("name").(string),
		},
	}

	ipblock, err := client.ReserveIPBlock(*ipblock)

	if err != nil {
		return fmt.Errorf("An error occured while reserving an ip block: %s", err)
	}
	d.SetId(ipblock.ID)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, ipblock.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceProfitBricksIPBlockRead(d, meta)
}

func resourceProfitBricksIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	ipblock, err := client.GetIPBlock(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching an ip block ID %s %s", d.Id(), err)
	}

	log.Printf("[INFO] IPS: %s", strings.Join(ipblock.Properties.IPs, ","))

	d.Set("ips", ipblock.Properties.IPs)
	d.Set("location", ipblock.Properties.Location)
	d.Set("size", ipblock.Properties.Size)
	d.Set("name", ipblock.Properties.Name)

	return nil
}
func resourceProfitBricksIPBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	request := profitbricks.IPBlockProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		request.Name = n.(string)
	}

	_, err := client.UpdateIPBlock(d.Id(), request)

	if err != nil {
		return fmt.Errorf("An error occured while updating an ip block ID %s %s", d.Id(), err)
	}

	return nil

}

func resourceProfitBricksIPBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	resp, err := client.ReleaseIPBlock(d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while releasing an ipblock ID: %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
