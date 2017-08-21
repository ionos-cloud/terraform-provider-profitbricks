package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
	"time"
)

func resourceProfitBricksLanIPFailover() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksLanIPFailoverCreate,
		Read:   resourceProfitBricksLanIPFailoverRead,
		Update: resourceProfitBricksLanIPFailoverUpdate,
		Delete: resourceProfitBricksLanIPFailoverDelete,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nicuuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lan_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProfitBricksLanIPFailoverCreate(d *schema.ResourceData, meta interface{}) error {
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)
	if lanid == "" {
		return fmt.Errorf("'lan_id' is missing, please provide a valid lan ID ")
	}
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)
	properties := &profitbricks.LanProperties{}

	properties.IpFailover = []profitbricks.IpFailover{}
	properties.IpFailover = append(properties.IpFailover, profitbricks.IpFailover{
		Ip:      ip,
		NicUuid: nicUuid,
	})

	if properties != nil {
		lan := profitbricks.PatchLan(dcid, lanid, *properties)
		if lan.StatusCode > 299 {
			return fmt.Errorf("An error occured while patching a lans failover group  %s %s", lanid, lan.Response)
		}
		err := waitTillProvisioned(meta, lan.Headers.Get("Location"))
		if err != nil {
			return err
		}
		d.SetId(lan.Id)
	}
	return resourceProfitBricksLanIPFailoverRead(d, meta)
}

func resourceProfitBricksLanIPFailoverRead(d *schema.ResourceData, meta interface{}) error {
	lan := profitbricks.GetLan(d.Get("datacenter_id").(string), d.Id())

	if lan.StatusCode > 299 {
		if lan.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a lan ID %s %s", d.Id(), lan.Response)
	}

	d.Set("public", lan.Properties.Public)
	d.Set("name", lan.Properties.Name)
	d.Set("ip_failover", lan.Properties.IpFailover)
	d.Set("datacenter_id", d.Get("datacenter_id").(string))
	return nil
}

func resourceProfitBricksLanIPFailoverUpdate(d *schema.ResourceData, meta interface{}) error {
	properties := &profitbricks.LanProperties{}
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)

	properties.IpFailover = []profitbricks.IpFailover{}
	properties.IpFailover = append(properties.IpFailover, profitbricks.IpFailover{
		Ip:      ip,
		NicUuid: nicUuid,
	})

	if properties != nil {
		lan := profitbricks.PatchLan(dcid, lanid, *properties)
		if lan.StatusCode > 299 {
			return fmt.Errorf("An error occured while patching a lan ID %s %s", d.Id(), lan.Response)
		}
		err := waitTillProvisioned(meta, lan.Headers.Get("Location"))
		if err != nil {
			return err
		}
	}
	return resourceProfitBricksLanIPFailoverRead(d, meta)
}

func resourceProfitBricksLanIPFailoverDelete(d *schema.ResourceData, meta interface{}) error {
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)

	properties := &profitbricks.LanProperties{}
	properties.IpFailover = make([]profitbricks.IpFailover, 0)

	resp := profitbricks.PatchLan(dcid, lanid, *properties)
	if resp.StatusCode > 299 {
		//try again in 50 seconds
		time.Sleep(50 * time.Second)
		resp = profitbricks.PatchLan(dcid, lanid, *properties)
		if resp.StatusCode > 299 && resp.StatusCode != 404 {
			return fmt.Errorf("An error occured while removing a lans ipfailover groups dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), string(resp.Response))
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
