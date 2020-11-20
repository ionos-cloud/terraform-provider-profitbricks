package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccProfitBricksPrivateCrossConnect_Basic(t *testing.T) {
	var privateCrossConnect profitbricks.PrivateCrossConnect
	privateCrossConnectName := "example"
	privateCrossConnectDescription := "example-description"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksprivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksprivateCrossConnectConfigBasic, privateCrossConnectName, privateCrossConnectDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksprivateCrossConnectExists("profitbricks_private_crossconnect.example", &privateCrossConnect),
					resource.TestCheckResourceAttr("profitbricks_private_crossconnect.example", "name", privateCrossConnectName),
					resource.TestCheckResourceAttr("profitbricks_private_crossconnect.example", "description", "example-description"),
				),
			},
			{
				Config: testAccCheckProfitBricksprivateCrossConnectConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksprivateCrossConnectExists("profitbricks_private_crossconnect.example", &privateCrossConnect),
					resource.TestCheckResourceAttr("profitbricks_private_crossconnect.example", "name", "example-renamed"),
					resource.TestCheckResourceAttr("profitbricks_private_crossconnect.example", "description", "example-description-updated"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksprivateCrossConnectDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_private_crossconnect" {
			continue
		}

		_, err := client.GetPrivateCrossConnect(rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("private cross-connect exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch private cross-connect %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckProfitBricksprivateCrossConnectExists(n string, privateCrossConnect *profitbricks.PrivateCrossConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
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

const testAccCheckProfitBricksprivateCrossConnectConfigBasic = `
resource "profitbricks_private_crossconnect" "example" {
  name        = "%s"
  description = "%s"
}`

const testAccCheckProfitBricksprivateCrossConnectConfigUpdate = `
resource "profitbricks_private_crossconnect" "example" {
  name        = "example-renamed"
  description = "example-description-updated"
}`
