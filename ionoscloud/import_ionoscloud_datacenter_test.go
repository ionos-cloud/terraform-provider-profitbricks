package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIonosCloudDataCenter_ImportBasic(t *testing.T) {
	resourceName := "datacenter-importtest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosCloudDatacenterConfig_basic, resourceName),
			},

			{
				ResourceName:      fmt.Sprintf("ionoscloud_datacenter.foobar"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
