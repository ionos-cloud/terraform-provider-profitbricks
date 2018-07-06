package profitbricks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProfitBricksResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{resource}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceProfitBricksFirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{nic}/{firewall}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("server_id", parts[1])
	d.Set("nic_id", parts[2])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}

