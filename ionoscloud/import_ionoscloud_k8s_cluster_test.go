package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIonosCloudk8sCluster_ImportBasic(t *testing.T) {
	resourceName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudk8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosCloudk8sClusterConfigBasic, resourceName),
			},
			{
				ResourceName:            fmt.Sprintf("ionoscloud_k8s_cluster.%s", resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"maintenance_window.0.time"},
			},
		},
	})
}
