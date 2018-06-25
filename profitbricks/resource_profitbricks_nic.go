package profitbricks

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func resourceProfitBricksNic() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksNicCreate,
		Read:   resourceProfitBricksNicRead,
		Update: resourceProfitBricksNicUpdate,
		Delete: resourceProfitBricksNicDelete,
		Schema: map[string]*schema.Schema{

			"lan": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"firewall_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nat": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksNicCreate(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	nic := profitbricks.Nic{
		Properties: &profitbricks.NicProperties{
			Lan: d.Get("lan").(int),
		},
	}
	if _, ok := d.GetOk("name"); ok {
		nic.Properties.Name = d.Get("name").(string)
	}
	if _, ok := d.GetOkExists("dhcp"); ok {
		val := d.Get("dhcp").(bool)
		nic.Properties.Dhcp = &val
	}

	if _, ok := d.GetOk("ip"); ok {
		raw := d.Get("ip").(string)
		ips := strings.Split(raw, ",")
		nic.Properties.Ips = ips
	}
	if _, ok := d.GetOk("firewall_active"); ok {
		raw := d.Get("firewall_active").(bool)
		nic.Properties.FirewallActive = raw
	}
	if _, ok := d.GetOk("nat"); ok {
		raw := d.Get("nat").(bool)
		nic.Properties.Nat = raw
	}

	resp, err := connection.CreateNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), nic)
	if err != nil {
		return fmt.Errorf("Error occured while creating a nic: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId(nic.ID)
	return resourceProfitBricksNicRead(d, meta)
}

func resourceProfitBricksNicRead(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	nic, err := connection.GetNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id())
	if err != nil {
		if err2, ok := err.(profitbricks.ApiError); ok {
			if err2.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a nic ID %s %s", d.Id(), err)
	}
	log.Printf("[INFO] LAN ON NIC: %d", nic.Properties.Lan)
	d.Set("dhcp", nic.Properties.Dhcp)
	d.Set("lan", nic.Properties.Lan)
	d.Set("name", nic.Properties.Name)
	d.Set("ips", nic.Properties.Ips)

	return nil
}

func resourceProfitBricksNicUpdate(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	properties := profitbricks.NicProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")

		properties.Name = n.(string)
	}
	if d.HasChange("lan") {
		_, n := d.GetChange("lan")
		properties.Lan = n.(int)
	}
	n := d.Get("dhcp").(bool)
	properties.Dhcp = &n

	if d.HasChange("ip") {
		_, raw := d.GetChange("ip")
		ips := strings.Split(raw.(string), ",")
		properties.Ips = ips
	}
	if d.HasChange("nat") {
		_, raw := d.GetChange("nat")
		nat := raw.(bool)
		properties.Nat = nat
	}

	nic, err := connection.UpdateNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id(), properties)

	if err != nil {
		return fmt.Errorf("Error occured while updating a nic: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, nic.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksNicRead(d, meta)
}

func resourceProfitBricksNicDelete(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	resp, err := connection.DeleteNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("An error occured while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
