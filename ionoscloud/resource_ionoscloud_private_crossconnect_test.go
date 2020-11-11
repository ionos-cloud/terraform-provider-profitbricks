package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudPrivateCrossConnect_Basic(t *testing.T) {
	var privateCrossConnect ionoscloud.PrivateCrossConnect
	privateCrossConnectName := "example"
	privateCrossConnectDescription := "example-description"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudprivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosCloudprivateCrossConnectConfigBasic, privateCrossConnectName, privateCrossConnectDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudprivateCrossConnectExists("ionoscloud_private_crossconnect.example", &privateCrossConnect),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "name", privateCrossConnectName),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "description", "example-description"),
				),
			},
			{
				Config: testAccCheckIonosCloudprivateCrossConnectConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudprivateCrossConnectExists("ionoscloud_private_crossconnect.example", &privateCrossConnect),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "name", "example-renamed"),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "description", "example-description-updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudprivateCrossConnectDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_private_crossconnect" {
			continue
		}

		_, err := client.GetPrivateCrossConnect(rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("private cross-connect exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch private cross-connect %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudprivateCrossConnectExists(n string, privateCrossConnect *ionoscloud.PrivateCrossConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundPrivateCrossConnect, err := client.GetPrivateCrossConnect(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching private cross-connect: %s", rs.Primary.ID)
		}
		if foundPrivateCrossConnect.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		privateCrossConnect = foundPrivateCrossConnect

		return nil
	}
}

const testAccCheckIonosCloudprivateCrossConnectConfigBasic = `
resource "ionoscloud_private_crossconnect" "example" {
  name        = "%s"
  description = "%s"
}`

const testAccCheckIonosCloudprivateCrossConnectConfigUpdate = `
resource "ionoscloud_private_crossconnect" "example" {
  name        = "example-renamed"
  description = "example-description-updated"
}`
