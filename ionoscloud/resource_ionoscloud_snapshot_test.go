package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudSnapshot_Basic(t *testing.T) {
	var snapshot ionoscloud.Snapshot
	snapshotName := "terraform_snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksSnapshotConfig_basic, snapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudSnapshotExists("ionoscloud_snapshot.test_snapshot", &snapshot),
					resource.TestCheckResourceAttr("ionoscloud_snapshot.test_snapshot", "name", snapshotName),
				),
			},
			{
				Config: testAccCheckProfitbricksSnapshotConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_snapshot.test_snapshot", "name", snapshotName),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudSnapshotDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_snapshot" {
			continue
		}

		_, err := client.GetSnapshot(rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("Snapshot still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching Snapshot %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudSnapshotExists(n string, snapshot *ionoscloud.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudSnapshotExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundServer, err := client.GetSnapshot(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Snapshot: %s", rs.Primary.ID)
		}
		if foundServer.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		snapshot = foundServer

		return nil
	}
}

const testAccCheckProfitbricksSnapshotConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "snapshot-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "debian:9"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 2
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  volume_id = "${ionoscloud_server.webserver.boot_volume}"
  name = "%s"
}
`

const testAccCheckProfitbricksSnapshotConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "snapshot-test"
	location = "us/las"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "debian:9"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 2
    disk_type = "HDD"
}
  nic {
    lan = "${ionoscloud_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  volume_id = "${ionoscloud_server.webserver.boot_volume}"
  name = "terraform_snapshot"
}`
