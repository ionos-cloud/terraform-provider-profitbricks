package profitbricks

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccProfitBricksFirewall_ImportBasic(t *testing.T) {
	firewallName := "firewall"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksFirewallDestroyCheck,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckProfitbricksFirewallConfig_basic, firewallName),
			},

			resource.TestStep{
				ResourceName:      "profitbricks_firewall.webserver_http",
				ImportStateIdFunc: testAccProfitBricksFirewallImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProfitBricksFirewallImportStateId(s *terraform.State) (string, error){
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_firewall" {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
