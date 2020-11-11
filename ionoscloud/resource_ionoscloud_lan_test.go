package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudLan_Basic(t *testing.T) {
	var lan ionoscloud.Lan
	lanName := "lanName"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudLanDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksLanConfig_basic, lanName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudLanExists("ionoscloud_lan.webserver_lan", &lan),
					testAccCheckIonosCloudLanAttributes("ionoscloud_lan.webserver_lan", lanName),
					resource.TestCheckResourceAttr("ionoscloud_lan.webserver_lan", "name", lanName),
				),
			},
			{
				Config: testAccCheckProfitbricksLanConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudLanAttributes("ionoscloud_lan.webserver_lan", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_lan.webserver_lan", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudLanDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		_, err := client.GetLan(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("LAN still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching LAN %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudLanAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckIonosCloudLanExists(n string, lan *ionoscloud.Lan) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudLanExists: Not found: %s", n)
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
resource "ionoscloud_datacenter" "foobar" {
	name       = "lan-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "%s"
}`

const testAccCheckProfitbricksLanConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "lan-test"
	location = "us/las"
}
resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "updated"
}`
