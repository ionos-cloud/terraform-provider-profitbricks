package profitbricks

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func resourceProfitBricksIPBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksIPBlockCreate,
		Read:   resourceProfitBricksIPBlockRead,
		Delete: resourceProfitBricksIPBlockDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	connection := meta.(*profitbricks.Client)
	ipblock := profitbricks.IPBlock{
		Properties: profitbricks.IPBlockProperties{
			Size:     d.Get("size").(int),
			Location: d.Get("location").(string),
			Name:     d.Get("name").(string),
		},
	}

	resp, err := connection.ReserveIPBlock(ipblock)

	if err != nil {
		return fmt.Errorf("An error occured while reserving an ip block: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId(resp.ID)

	return resourceProfitBricksIPBlockRead(d, meta)
}

func resourceProfitBricksIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	ipblock, err := connection.GetIPBlock(d.Id())

	if err != nil {
		if err2, ok := err.(profitbricks.ApiError); ok {
			if err2.HttpStatusCode() == 404 {
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

func resourceProfitBricksIPBlockDelete(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	resp, err := connection.ReleaseIPBlock(d.Id())
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
