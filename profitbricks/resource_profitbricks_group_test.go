package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func TestAccProfitBricksGroup_Basic(t *testing.T) {
	var group profitbricks.Group
	groupName := "terraform test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksGroupConfig_basic, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksGroupExists("profitbricks_group.group", &group),
					testAccCheckProfitBricksGroupAttributes("profitbricks_group.group", groupName),
					resource.TestCheckResourceAttr("profitbricks_group.group", "name", groupName),
				),
			},
			{
				Config: testAccCheckProfitbricksGroupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksGroupAttributes("profitbricks_group.group", "updated"),
					resource.TestCheckResourceAttr("profitbricks_group.group", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksGroupDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		resp := profitbricks.GetGroup(rs.Primary.ID)

		if resp.StatusCode < 299 {
			return fmt.Errorf("group still exists %s %s", rs.Primary.ID, resp.Response)
		}
	}

	return nil
}

func testAccCheckProfitBricksGroupAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksGroupAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckProfitBricksGroupExists(n string, group *profitbricks.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundgroup := profitbricks.GetGroup(rs.Primary.ID)

		if foundgroup.StatusCode != 200 {
			return fmt.Errorf("Error occured while fetching Group: %s", rs.Primary.ID)
		}
		if foundgroup.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		group = &foundgroup

		return nil
	}
}

const testAccCheckProfitbricksGroupConfig_basic = `
resource "profitbricks_group" "group" {
  name = "%s"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}`

const testAccCheckProfitbricksGroupConfig_update = `
resource "profitbricks_group" "group" {
  name = "updated"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
}
`
