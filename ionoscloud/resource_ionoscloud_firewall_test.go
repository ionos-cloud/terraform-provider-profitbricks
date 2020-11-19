package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudFirewall_Basic(t *testing.T) {
	var firewall ionoscloud.FirewallRule
	firewallName := "firewall"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksFirewallConfig_basic, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudFirewallExists("ionoscloud_firewall.webserver_http", &firewall),
					testAccCheckIonosCloudFirewallAttributes("ionoscloud_firewall.webserver_http", firewallName),
					resource.TestCheckResourceAttr("ionoscloud_firewall.webserver_http", "name", firewallName),
				),
			},
			{
				Config: testAccCheckProfitbricksFirewallConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudFirewallAttributes("ionoscloud_firewall.webserver_http", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_firewall.webserver_http", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudFirewallDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_firewall" {
			continue
		}

		_, err := client.GetFirewallRule(rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("Firewall still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching Firewall %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudFirewallAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudFirewallAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckIonosCloudFirewallExists(n string, firewall *ionoscloud.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckIonosCloudFirewallExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundServer, err := client.GetFirewallRule(rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Firewall rule: %s", rs.Primary.ID)
		}
		if foundServer.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		firewall = foundServer

		return nil
	}
}

const testAccCheckProfitbricksFirewallConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "firewall-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu-16.04"
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
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "TCP"
  name = "%s"
  port_range_start = 80
  port_range_end = 80
}`

const testAccCheckProfitbricksFirewallConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "firewall-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu-16.04"
	image_password = "test1234"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
}
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "TCP"
  name = "updated"
  port_range_start = 80
  port_range_end = 80
}`
