package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccProfitBricksLan_Basic(t *testing.T) {
	var lan profitbricks.Lan
	lanName := "lanName"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksLanDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksLanConfig_basic, lanName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksLanExists("profitbricks_lan.webserver_lan", &lan),
					testAccCheckProfitBricksLanAttributes("profitbricks_lan.webserver_lan", lanName),
					resource.TestCheckResourceAttr("profitbricks_lan.webserver_lan", "name", lanName),
				),
			},
			{
				Config: testAccCheckProfitbricksLanConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksLanAttributes("profitbricks_lan.webserver_lan", "updated"),
					resource.TestCheckResourceAttr("profitbricks_lan.webserver_lan", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksLanDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_datacenter" {
			continue
		}

		_, err := client.GetLan(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("LAN still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching LAN %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckProfitBricksLanAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckProfitBricksLanExists(n string, lan *profitbricks.Lan) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksLanExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundLan, err := client.GetLan(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Server: %s", rs.Primary.ID)
		}
		if foundLan.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		lan = foundLan

		return nil
	}
}

const testAccCheckProfitbricksLanConfig_basic = `
resource "profitbricks_datacenter" "foobar" {
	name       = "lan-test"
	location = "us/las"
}

resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "%s"
}`

const testAccCheckProfitbricksLanConfig_update = `
resource "profitbricks_datacenter" "foobar" {
	name       = "lan-test"
	location = "us/las"
}
resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "updated"
}`
