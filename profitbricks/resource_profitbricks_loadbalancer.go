package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func resourceProfitBricksLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksLoadbalancerCreate,
		Read:   resourceProfitBricksLoadbalancerRead,
		Update: resourceProfitBricksLoadbalancerUpdate,
		Delete: resourceProfitBricksLoadbalancerDelete,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"datacenter_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"nic_ids": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nic_id": &schema.Schema{
				Type:     schema.TypeList,
				Removed:  "Use nic_ids instead",
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceProfitBricksLoadbalancerCreate(d *schema.ResourceData, meta interface{}) error {

	raw_ids := d.Get("nic_ids").([]interface{})
	nic_ids := []profitbricks.Nic{}

	for _, id := range raw_ids {
		nic_ids = append(nic_ids, profitbricks.Nic{Id: id.(string)})
	}

	lb := profitbricks.Loadbalancer{
		Properties: profitbricks.LoadbalancerProperties{
			Name: d.Get("name").(string),
		},
		Entities: profitbricks.LoadbalancerEntities{
			Balancednics: &profitbricks.BalancedNics{
				Items: nic_ids,
			},
		},
	}

	lb = profitbricks.CreateLoadbalancer(d.Get("datacenter_id").(string), lb)

	if lb.StatusCode > 299 {
		return fmt.Errorf("Error occured while creating a loadbalancer %s", lb.Response)
	}
	err := waitTillProvisioned(meta, lb.Headers.Get("Location"))

	if err != nil {
		return err
	}

	d.SetId(lb.Id)

	return resourceProfitBricksLoadbalancerRead(d, meta)
}

func resourceProfitBricksLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	lb := profitbricks.GetLoadbalancer(d.Get("datacenter_id").(string), d.Id())

	if lb.StatusCode > 299 {
		if lb.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a lan ID %s %s", d.Id(), lb.Response)
	}

	d.Set("name", lb.Properties.Name)
	d.Set("ip", lb.Properties.Ip)
	d.Set("dhcp", lb.Properties.Dhcp)

	return nil
}

func resourceProfitBricksLoadbalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	properties := profitbricks.LoadbalancerProperties{}
	if d.HasChange("name") {
		_, new := d.GetChange("name")
		properties.Name = new.(string)
	}
	if d.HasChange("ip") {
		_, new := d.GetChange("ip")
		properties.Ip = new.(string)
	}
	if d.HasChange("dhcp") {
		_, new := d.GetChange("dhcp")
		properties.Dhcp = new.(bool)
	}

	if d.HasChange("nic_ids") {
		old, new := d.GetChange("nic_ids")

		oldList := old.([]interface{})

		for _, o := range oldList {

			resp := profitbricks.DeleteBalancedNic(d.Get("datacenter_id").(string), d.Id(), o.(string))
			if resp.StatusCode > 299 {
				return fmt.Errorf("Error occured while deleting a balanced nic: %s", string(resp.Body))
			}
			err := waitTillProvisioned(meta, resp.Headers.Get("Location"))
			if err != nil {
				return err
			}
		}

		newList := new.([]interface{})

		for _, o := range newList {
			nic := profitbricks.AssociateNic(d.Get("datacenter_id").(string), d.Id(), o.(string))
			if nic.StatusCode > 299 {
				return fmt.Errorf("Error occured while deleting a balanced nic: %s", nic.Response)
			}
			err := waitTillProvisioned(meta, nic.Headers.Get("Location"))
			if err != nil {
				return err
			}
		}

	}

	return resourceProfitBricksLoadbalancerRead(d, meta)
}

func resourceProfitBricksLoadbalancerDelete(d *schema.ResourceData, meta interface{}) error {
	resp := profitbricks.DeleteLoadbalancer(d.Get("datacenter_id").(string), d.Id())

	if resp.StatusCode > 299 {
		return fmt.Errorf("Error occured while deleting a loadbalancer: %s", string(resp.Body))
	}

	err := waitTillProvisioned(meta, resp.Headers.Get("Location"))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
