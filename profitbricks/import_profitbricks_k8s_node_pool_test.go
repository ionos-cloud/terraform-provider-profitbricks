package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
