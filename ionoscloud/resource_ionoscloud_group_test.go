package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudGroup_Basic(t *testing.T) {
	var group ionoscloud.Group
	groupName := "terraform test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksGroupConfig_basic, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudGroupExists("ionoscloud_group.group", &group),
					testAccCheckIonosCloudGroupAttributes("ionoscloud_group.group", groupName),
					resource.TestCheckResourceAttr("ionoscloud_group.group", "name", groupName),
				),
			},
			{
				Config: testAccCheckProfitbricksGroupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudGroupAttributes("ionoscloud_group.group", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_group.group", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		_, err := client.GetGroup(rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("group still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching Group %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudGroupAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudGroupAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckIonosCloudGroupExists(n string, group *ionoscloud.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundgroup, err := client.GetGroup(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Group: %s", rs.Primary.ID)
		}
		if foundgroup.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		group = foundgroup

		return nil
	}
}

const testAccCheckProfitbricksGroupConfig_basic = `
resource "ionoscloud_group" "group" {
  name = "%s"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}`

const testAccCheckProfitbricksGroupConfig_update = `
resource "ionoscloud_group" "group" {
  name = "updated"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
}
`
