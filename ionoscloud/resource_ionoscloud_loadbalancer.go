package ionoscloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceIonosCloudLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceIonosCloudLoadbalancerCreate,
		Read:   resourceIonosCloudLoadbalancerRead,
		Update: resourceIonosCloudLoadbalancerUpdate,
		Delete: resourceIonosCloudLoadbalancerDelete,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nic_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceIonosCloudLoadbalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	raw_ids := d.Get("nic_ids").([]interface{})
	nic_ids := []ionoscloud.Nic{}

	for _, id := range raw_ids {
		nic_ids = append(nic_ids, ionoscloud.Nic{ID: id.(string)})
	}

	lb := &ionoscloud.Loadbalancer{
		Properties: ionoscloud.LoadbalancerProperties{
			Name: d.Get("name").(string),
		},
		Entities: ionoscloud.LoadbalancerEntities{
			Balancednics: &ionoscloud.BalancedNics{
				Items: nic_ids,
			},
		},
	}

	lb, err := client.CreateLoadbalancer(d.Get("datacenter_id").(string), *lb)

	if err != nil {
		return fmt.Errorf("Error occured while creating a loadbalancer %s", err)
	}
	d.SetId(lb.ID)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, lb.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceIonosCloudLoadbalancerRead(d, meta)
}

func resourceIonosCloudLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	lb, err := client.GetLoadbalancer(d.Get("datacenter_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a lan ID %s %s", d.Id(), err)
	}

	d.Set("name", lb.Properties.Name)
	d.Set("ip", lb.Properties.IP)
	d.Set("dhcp", lb.Properties.Dhcp)

	return nil
}

func resourceIonosCloudLoadbalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	properties := ionoscloud.LoadbalancerProperties{}
	if d.HasChange("name") {
		_, new := d.GetChange("name")
		properties.Name = new.(string)
	}
	if d.HasChange("ip") {
		_, new := d.GetChange("ip")
		properties.IP = new.(string)
	}
	if d.HasChange("dhcp") {
		_, new := d.GetChange("dhcp")
		properties.Dhcp = new.(bool)
	}

	if d.HasChange("nic_ids") {
		old, new := d.GetChange("nic_ids")

		oldList := old.([]interface{})

		for _, o := range oldList {

			resp, err := client.DeleteBalancedNic(d.Get("datacenter_id").(string), d.Id(), o.(string))
			if err != nil {
				return fmt.Errorf("Error occured while deleting a balanced nic: %s", err)
			}

			// Wait, catching any errors
			_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutUpdate).WaitForState()
			if errState != nil {
				return errState
			}
		}

		newList := new.([]interface{})

		for _, o := range newList {
			nic, err := client.AssociateNic(d.Get("datacenter_id").(string), d.Id(), o.(string))
			if err != nil {
				return fmt.Errorf("Error occured while deleting a balanced nic: %s", err)
			}

			// Wait, catching any errors
			_, errState := getStateChangeConf(meta, d, nic.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
			if errState != nil {
				return errState
			}

		}

	}

	return resourceIonosCloudLoadbalancerRead(d, meta)
}

func resourceIonosCloudLoadbalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.Client)
	resp, err := client.DeleteLoadbalancer(d.Get("datacenter_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error occured while deleting a loadbalancer: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
