package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func TestAccProfitBricksShare_Basic(t *testing.T) {
	var share profitbricks.Share
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksShareDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksShareConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksShareExists("profitbricks_share.share", &share),
					resource.TestCheckResourceAttr("profitbricks_share.share", "share_privilege", "true"),
				),
			},
			{
				Config: testAccCheckProfitbricksShareConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("profitbricks_share.share", "share_privilege", "false"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksShareDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		resp := profitbricks.GetShare(rs.Primary.Attributes["group_id"], rs.Primary.Attributes["resource_id"])

		if resp.StatusCode < 299 {
			return fmt.Errorf("share for resource %s still exists in group %s %s", rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["group_id"], resp.Response)
		}
	}

	return nil
}

func testAccCheckProfitBricksShareExists(n string, share *profitbricks.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksShareExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		grp_id := rs.Primary.Attributes["group_id"]
		resource_id := rs.Primary.Attributes["resource_id"]

		foundshare := profitbricks.GetShare(grp_id, resource_id)

		if foundshare.StatusCode != 200 {
			return fmt.Errorf("Error occured while fetching Share of resource  %s in group %s", rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["group_id"])
		}
		if foundshare.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		share = &foundshare

		return nil
	}
}

const testAccCheckProfitbricksShareConfig_basic = `
resource "profitbricks_datacenter" "foobar" {
	name       = "terraform test"
	location = "us/las"
}

resource "profitbricks_group" "group" {
  name = "terraform test"
  create_dataCenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "profitbricks_share" "share" {
  group_id = "${profitbricks_group.group.id}"
  resource_id = "${profitbricks_datacenter.foobar.id}"
  edit_privilege = true
  share_privilege = true
}`

const testAccCheckProfitbricksShareConfig_update = `
resource "profitbricks_datacenter" "foobar" {
	name       = "terraform test"
	location = "us/las"
}

resource "profitbricks_group" "group" {
  name = "terraform test"
  create_dataCenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "profitbricks_share" "share" {
group_id = "${profitbricks_group.group.id}"
  resource_id = "${profitbricks_datacenter.foobar.id}"
  edit_privilege = true
  share_privilege = false
}
`
