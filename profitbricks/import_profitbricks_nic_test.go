package profitbricks

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccProfitBricksNic_ImportBasic(t *testing.T) {
	volumeName := "volume"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksNicDestroyCheck,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckProfitbricksNicConfig_basic, volumeName),
			},

			resource.TestStep{
				ResourceName:      "profitbricks_nic.database_nic",
				ImportStateIdFunc: testAccProfitBricksNicImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProfitBricksNicImportStateId(s *terraform.State) (string, error){
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_nic" && rs.Primary.ID != "database_nic"{
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
