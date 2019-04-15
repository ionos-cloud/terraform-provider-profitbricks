package profitbricks

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func resourceProfitBricksLan() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksLanCreate,
		Read:   resourceProfitBricksLanRead,
		Update: resourceProfitBricksLanUpdate,
		Delete: resourceProfitBricksLanDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProfitBricksResourceImport,
		},
		Schema: map[string]*schema.Schema{

			"public": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksLanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	request := profitbricks.Lan{
		Properties: profitbricks.LanProperties{
			Public: d.Get("public").(bool),
		},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("name"))
	if d.Get("name") != nil {
		request.Properties.Name = d.Get("name").(string)
	}

	lan, err := client.CreateLan(d.Get("datacenter_id").(string), request)

	log.Printf("[DEBUG] LAN ID: %s", lan.ID)
	log.Printf("[DEBUG] LAN RESPONSE: %s", lan.Response)

	if err != nil {
		return fmt.Errorf("An error occured while creating a lan: %s", err)
	}

	d.SetId(lan.ID)
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, lan.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceProfitBricksLanRead(d, meta)
}

func resourceProfitBricksLanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	lan, err := client.GetLan(d.Get("datacenter_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a lan ID %s %s", d.Id(), err)
	}

	d.Set("public", lan.Properties.Public)
	d.Set("name", lan.Properties.Name)
	d.Set("ip_failover", lan.Properties.IPFailover)
	d.Set("datacenter_id", d.Get("datacenter_id").(string))
	return nil
}

func resourceProfitBricksLanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	properties := &profitbricks.LanProperties{}
	newValue := d.Get("public")
	properties.Public = newValue.(bool)
	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		properties.Name = newValue.(string)
	}

	if properties != nil {
		lan, err := client.UpdateLan(d.Get("datacenter_id").(string), d.Id(), *properties)
		if err != nil {
			return fmt.Errorf("An error occured while patching a lan ID %s %s", d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, lan.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}

	}
	return resourceProfitBricksLanRead(d, meta)
}

func resourceProfitBricksLanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)
	resp, err := client.DeleteLan(dcId, d.Id())
	if err != nil {
		//try again in 120 seconds
		time.Sleep(120 * time.Second)
		resp, err = client.DeleteLan(dcId, d.Id())

		if err != nil {
			if apiError, ok := err.(profitbricks.ApiError); ok {
				if apiError.HttpStatusCode() != 404 {
					return fmt.Errorf("An error occured while deleting a lan dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
				}
			}
		}
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
