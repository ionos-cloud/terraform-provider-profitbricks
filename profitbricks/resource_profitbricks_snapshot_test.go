package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
	"testing"
)

func TestAccProfitBricksSnapshot_Basic(t *testing.T) {
	var snapshot profitbricks.Snapshot
	snapshotName := "terraform_snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksSnapshotConfig_basic, snapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksSnapshotExists("profitbricks_snapshot.test_snapshot", &snapshot),
					resource.TestCheckResourceAttr("profitbricks_snapshot.test_snapshot", "name", snapshotName),
				),
			},
			{
				Config: testAccCheckProfitbricksSnapshotConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("profitbricks_snapshot.test_snapshot", "name", snapshotName),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksSnapshotDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_snapshot" {
			continue
		}

		resp := profitbricks.GetSnapshot(rs.Primary.ID)

		if resp.StatusCode < 299 {
			return fmt.Errorf("Snapshot still exists %s %s", rs.Primary.ID, resp.Response)
		}
	}

	return nil
}

func testAccCheckProfitBricksSnapshotExists(n string, snapshot *profitbricks.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksSnapshotExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundServer := profitbricks.GetSnapshot(rs.Primary.ID)

		if foundServer.StatusCode != 200 {
			return fmt.Errorf("Error occured while fetching Snapshot: %s", rs.Primary.ID)
		}
		if foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		snapshot = &foundServer

		return nil
	}
}

const testAccCheckProfitbricksSnapshotConfig_basic = `
resource "profitbricks_datacenter" "foobar" {
	name       = "snapshot-test"
	location = "us/las"
}

resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "profitbricks_server" "webserver" {
  name = "webserver"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
    name = "system"
    size = 2
    disk_type = "HDD"
    image_name = "debian:9"
    image_password = "K3tTj8G14a3EgKyNeeiY"
}
  nic {
    lan = "${profitbricks_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "profitbricks_snapshot" "test_snapshot" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  volume_id = "${profitbricks_server.webserver.boot_volume}"
  name = "%s"
}
`

const testAccCheckProfitbricksSnapshotConfig_update = `
resource "profitbricks_datacenter" "foobar" {
	name       = "snapshot-test"
	location = "us/las"
}

resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "profitbricks_server" "webserver" {
  name = "webserver"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
    name = "system"
    size = 2
    disk_type = "HDD"
    image_name = "debian:9"
    image_password = "K3tTj8G14a3EgKyNeeiY"
}
  nic {
    lan = "${profitbricks_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "profitbricks_snapshot" "test_snapshot" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  volume_id = "${profitbricks_server.webserver.boot_volume}"
  name = "terraform_snapshot"
}`
