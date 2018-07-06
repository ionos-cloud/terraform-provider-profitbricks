package profitbricks

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccProfitBricksLan_ImportBasic(t *testing.T) {
	lanName := "lanName"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksLanDestroyCheck,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckProfitbricksLanConfig_basic, lanName),
			},

			resource.TestStep{
				ResourceName:      "profitbricks_lan.webserver_lan",
				ImportStateIdFunc: testAccProfitBricksLanImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProfitBricksLanImportStateId(s *terraform.State) (string, error){
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_lan" {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
