package profitbricks

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccProfitBricksServer_ImportBasic(t *testing.T) {
	resourceName := "server-importtest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksServerDestroyCheck,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckProfitbricksServerConfig_basic, resourceName),
			},

			resource.TestStep{
				ResourceName:            "profitbricks_server.webserver",
				ImportStateIdFunc:       testAccProfitBricksServerImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_password", "boot_image"},
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

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"], rs.Primary.Attributes["primary_nic"])
	}

	return importID, nil
}
