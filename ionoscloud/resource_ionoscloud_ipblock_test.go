package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudIPBlock_Basic(t *testing.T) {
	var ipblock ionoscloud.IPBlock
	location := "us/las"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksIPBlockConfig_basic, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudIPBlockExists("ionoscloud_ipblock.webserver_ip", &ipblock),
					testAccCheckIonosCloudIPBlockAttributes("ionoscloud_ipblock.webserver_ip", location),
					resource.TestCheckResourceAttr("ionoscloud_ipblock.webserver_ip", "location", location),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksIPBlockConfig_update, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudIPBlockExists("ionoscloud_ipblock.webserver_ip", &ipblock),
					testAccCheckIonosCloudIPBlockAttributes("ionoscloud_ipblock.webserver_ip", location),
					resource.TestCheckResourceAttr("ionoscloud_ipblock.webserver_ip", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudIPBlockDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_ipblock" {
			continue
		}

		_, err := client.GetIPBlock(rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("IPBlock still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching IPBlock %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudIPBlockAttributes(n string, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["location"] != location {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["location"])
		}

		return nil
	}
}

func testAccCheckIonosCloudIPBlockExists(n string, ipblock *ionoscloud.IPBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudIPBlockExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundIP, err := client.GetIPBlock(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching IP Block: %s", rs.Primary.ID)
		}
		if foundIP.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		ipblock = foundIP

		return nil
	}
}

const testAccCheckProfitbricksIPBlockConfig_basic = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "ipblock TF test"
}`

const testAccCheckProfitbricksIPBlockConfig_update = `
resource "ionoscloud_ipblock" "webserver_ip" {
  location = "%s"
  size = 1
  name = "updated"
}`
