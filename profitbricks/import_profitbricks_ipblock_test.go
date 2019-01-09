package profitbricks

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccProfitBricksIPBlock_ImportBasic(t *testing.T) {
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksIPBlockConfig_basic, location),
			},

			{
				ResourceName:      "profitbricks_ipblock.webserver_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
