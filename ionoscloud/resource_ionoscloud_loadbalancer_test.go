package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudLoadbalancer_Basic(t *testing.T) {
	var loadbalancer ionoscloud.Loadbalancer
	lbName := "loadbalancer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudLoadbalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksLoadbalancerConfig_basic, lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudLoadbalancerExists("ionoscloud_loadbalancer.example", &loadbalancer),
					testAccCheckIonosCloudLoadbalancerAttributes("ionoscloud_loadbalancer.example", lbName),
					resource.TestCheckResourceAttr("ionoscloud_loadbalancer.example", "name", lbName),
				),
			},
			{
				Config: testAccCheckProfitbricksLoadbalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudLoadbalancerAttributes("ionoscloud_loadbalancer.example", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_loadbalancer.example", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudLoadbalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_loadbalancer" {
			continue
		}

		_, err := client.GetLoadbalancer(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if err != nil {
			_, err := client.DeleteDatacenter(rs.Primary.Attributes["datacenter_id"])

			if apiError, ok := err.(ionoscloud.ApiError); ok {
				if apiError.HttpStatusCode() != 404 {
					return fmt.Errorf("loadbalancer still exists %s %s", rs.Primary.ID, apiError)
				}
			} else {
				return fmt.Errorf("Unable to fetching loadbalancer %s %s", rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccCheckIonosCloudLoadbalancerAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudLoadbalancerAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckIonosCloudLoadbalancerExists(n string, loadbalancer *ionoscloud.Loadbalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudLoadbalancerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundLB, err := client.GetLoadbalancer(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Loadbalancer: %s", rs.Primary.ID)
		}
		if foundLB.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		loadbalancer = foundLB

		return nil
	}
}

const testAccCheckProfitbricksLoadbalancerConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "loadbalancer-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "ubuntu-16.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
}
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = "2"
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_loadbalancer" "example" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  nic_ids = ["${ionoscloud_nic.database_nic.id}"]
  name = "%s"
  dhcp = true
}`

const testAccCheckProfitbricksLoadbalancerConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "loadbalancer-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "ubuntu-16.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_nic" "database_nic1" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = "2"
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_nic" "database_nic2" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = "3"
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_loadbalancer" "example" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  nic_ids = ["${ionoscloud_nic.database_nic1.id}","${ionoscloud_nic.database_nic2.id}"]
  name = "updated"
  dhcp = true
}`
