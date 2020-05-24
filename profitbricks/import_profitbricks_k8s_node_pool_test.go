package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccProfitBricksk8sNodepool_ImportBasic(t *testing.T) {
	resourceName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksk8sNodepoolConfigBasic, resourceName),
			},
			{
				ResourceName:      fmt.Sprintf("profitbricks_k8s_node_pool.%s", resourceName),
				ImportStateIdFunc: testAccProfitBricksk8sNodepoolImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProfitBricksk8sNodepoolImportStateID(s *terraform.State) (string, error) {
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_k8s_node_pool" {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)
	}

	return importID, nil
}
