package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIonosCloudIPBlock_ImportBasic(t *testing.T) {
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksIPBlockConfig_basic, location),
			},

			{
				ResourceName:      "ionoscloud_ipblock.webserver_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
