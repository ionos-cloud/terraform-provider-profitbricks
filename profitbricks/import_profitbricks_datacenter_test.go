package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccProfitBricksDataCenter_ImportBasic(t *testing.T) {
	resourceName := "datacenter-importtest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksDatacenterConfig_basic, resourceName),
			},

			{
				ResourceName:      fmt.Sprintf("profitbricks_datacenter.foobar"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
