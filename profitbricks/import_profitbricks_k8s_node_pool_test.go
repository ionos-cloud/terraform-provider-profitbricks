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
		CheckDestroy: testAccCheckDProfitBricksk8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksk8sClusterConfigBasic, resourceName),
			},
			{
				ResourceName:      fmt.Sprintf("profitbricks_k8s_cluster.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
