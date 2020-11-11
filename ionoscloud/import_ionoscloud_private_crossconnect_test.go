package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIonosCloudprivateCrossConnect_ImportBasic(t *testing.T) {
	resourceName := "example"
	resourceDescription := "example-description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudprivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosCloudprivateCrossConnectConfigBasic, resourceName, resourceDescription),
			},
			{
				ResourceName:      fmt.Sprintf("ionoscloud_private_crossconnect.%s", resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
