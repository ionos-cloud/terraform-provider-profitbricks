package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudDataCenter_Basic(t *testing.T) {
	var datacenter ionoscloud.Datacenter
	dc_name := "datacenter-test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosCloudDatacenterConfig_basic, dc_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudDatacenterExists("ionoscloud_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("ionoscloud_datacenter.foobar", "name", dc_name),
				),
			},
			{
				Config: testAccCheckIonosCloudDatacenterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudDatacenterExists("ionoscloud_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("ionoscloud_datacenter.foobar", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudDatacenterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		_, err := client.GetDatacenter(rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("DataCenter still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching DataCenter %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudDatacenterExists(n string, datacenter *ionoscloud.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundDC, err := client.GetDatacenter(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching DC: %s", rs.Primary.ID)
		}
		if foundDC.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		datacenter = foundDC

		return nil
	}
}

const testAccCheckIonosCloudDatacenterConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "%s"
	location = "us/las"
}`

const testAccCheckIonosCloudDatacenterConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       =  "updated"
	location = "us/las"
}`
