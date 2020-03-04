package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccProfitBricksServer_Basic(t *testing.T) {
	var server profitbricks.Server
	serverName := "webserver"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksServerConfig_basic, serverName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksServerExists("profitbricks_server.webserver", &server),
					testAccCheckProfitBricksServerAttributes("profitbricks_server.webserver", serverName),
					resource.TestCheckResourceAttr("profitbricks_server.webserver", "name", serverName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksServerConfig_basicdep, serverName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksServerExists("profitbricks_server.webserver", &server),
					testAccCheckProfitBricksServerAttributes("profitbricks_server.webserver", serverName),
					resource.TestCheckResourceAttr("profitbricks_server.webserver", "name", serverName),
				),
			},
			{
				Config: testAccCheckProfitbricksServerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksServerAttributes("profitbricks_server.webserver", "updated"),
					resource.TestCheckResourceAttr("profitbricks_server.webserver", "name", "updated"),
					resource.TestCheckResourceAttr("profitbricks_server.webserver", "nic.0.dhcp", "false"),
					resource.TestCheckResourceAttr("profitbricks_server.webserver", "nic.0.firewall_active", "false"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_datacenter" {
			continue
		}

		_, err := client.GetServer(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("Server still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching Server %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckProfitBricksServerAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksServerAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckProfitBricksServerExists(n string, server *profitbricks.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksServerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundServer, err := client.GetServer(rs.Primary.Attributes["datacenter_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching Server: %s", rs.Primary.ID)
		}
		if foundServer.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		server = foundServer

		return nil
	}
}

const testAccCheckProfitbricksServerConfig_basic = `
resource "profitbricks_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "profitbricks_server" "webserver" {
  name = "%s"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu:latest"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
}
  nic {
    lan = "${profitbricks_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckProfitbricksServerConfig_basicdep = `
resource "profitbricks_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "profitbricks_server" "webserver" {
  name = "%s"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
		image_name ="ubuntu:latest"
		image_password = "K3tTj8G14a3EgKyNeeiY"
    name = "system"
    size = 5
    disk_type = "SSD"
  }
  nic {
    lan = "${profitbricks_lan.webserver_lan.id}"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

const testAccCheckProfitbricksServerConfig_update = `
resource "profitbricks_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "profitbricks_lan" "webserver_lan" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "public"
}

resource "profitbricks_server" "webserver" {
  name = "updated"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "ubuntu:latest"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
}
  nic {
    lan = "${profitbricks_lan.webserver_lan.id}"
    dhcp = false
    firewall_active = false
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}`

func Test_Update(t *testing.T) {

}
