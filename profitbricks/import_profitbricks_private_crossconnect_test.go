package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccProfitBricksprivateCrossConnect_ImportBasic(t *testing.T) {
	resourceName := "example"
	resourceDescription := "example-description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksprivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksprivateCrossConnectConfigBasic, resourceName, resourceDescription),
			},
			{
				ResourceName:      fmt.Sprintf("profitbricks_private_crossconnect.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
