package profitbricks

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccProfitBricksServer_ImportBasic(t *testing.T) {
	resourceName := "server-importtest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksServerConfig_basic, resourceName),
			},

			{
				ResourceName:            "profitbricks_server.webserver",
				ImportStateIdFunc:       testAccProfitBricksServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "ssh_key_path.#", "image_name"},
			},
		},
	})
}

func testAccProfitBricksServerImportStateId(s *terraform.State) (string, error) {
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_server" {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"], rs.Primary.Attributes["primary_nic"], rs.Primary.Attributes["firewallrule_id"])
	}

	return importID, nil
}
