package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func TestAccProfitBricksLanIPFailover_Basic(t *testing.T) {
	var lan profitbricks.Lan
	var ipfailover profitbricks.IpFailover

	testDeleted := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[n]
			if ok {
				return fmt.Errorf("Expected Failover group  %s to be deleted.", n)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksLanIPFailoverDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksLanIPFailoverConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanIPFailoverGroupExists("profitbricks_ipfailover.failovertest", &lan, &ipfailover),
				),
			},
			{
				Config: testAccCheckProfitbricksLanIPFailoverConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("profitbricks_ipfailover.failovertest"),
				),
			},
		},
	})
}

func testAccCheckLanIPFailoverGroupExists(n string, lan *profitbricks.Lan, failover *profitbricks.IpFailover) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID` is set")
		}

		lanId := rs.Primary.Attributes["lan_id"]
		nicUuid := rs.Primary.Attributes["nicuuid"]

		lan := profitbricks.GetLan(rs.Primary.Attributes["datacenter_id"], lanId)
		if lan.StatusCode > 299 {
			return fmt.Errorf("Lan %s not found.", lanId)
		}

		if len(lan.Properties.IpFailover) == 0 {
			return fmt.Errorf("Lan %s has no failover groups.", lanId)
		}
		found := false
		for _, fo := range lan.Properties.IpFailover {
			if fo.NicUuid == nicUuid {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Expected NIC %s to be a part of a failover group, but not found in lans %s failover groups", nicUuid, lanId)
		}

		return nil
	}
}

func testAccCheckDProfitBricksLanIPFailoverDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_ipfailover" {
			continue
		}
		nicUuid := rs.Primary.Attributes["nicuuid"]
		resp := profitbricks.GetLan(rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["lan_id"])
		found := false
		for _, fo := range resp.Properties.IpFailover {
			if fo.NicUuid == nicUuid {
				found = true
			}
		}
		if found {
			return fmt.Errorf("IP failover group with nicId %s still exists %s %s", nicUuid, rs.Primary.ID, resp.Response)
		}
	}

	return nil
}

const testAccCheckProfitbricksLanIPFailoverConfig_basic = `
resource "profitbricks_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "profitbricks_ipblock" "webserver_ip" {
  location = "us/las"
  size = 1
  name = "failover test"
}

resource "profitbricks_lan" "webserver_lan1" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "terraform test"
}

resource "profitbricks_server" "webserver" {
  name = "server"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
    image_name="centos:latest"
    image_password = "K3tTj8G14a3EgKyNeeiY"
}
  nic {
    lan = "${profitbricks_lan.webserver_lan1.id}"
    dhcp = true
    firewall_active = true
     ip ="${profitbricks_ipblock.webserver_ip.ips[0]}"
  }
}
resource "profitbricks_ipfailover" "failovertest" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  lan_id="${profitbricks_lan.webserver_lan1.id}"
  ip ="${profitbricks_ipblock.webserver_ip.ips[0]}"
  nicuuid= "${profitbricks_server.webserver.primary_nic}"
}`

const testAccCheckProfitbricksLanIPFailoverConfig_update = `
resource "profitbricks_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "profitbricks_ipblock" "webserver_ip" {
  location = "us/las"
  size = 1
  name = "failover test"
}

resource "profitbricks_lan" "webserver_lan1" {
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  public = true
  name = "terraform test"
}

resource "profitbricks_server" "webserver" {
  name = "server"
  datacenter_id = "${profitbricks_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
    image_name="centos:latest"
    image_password = "K3tTj8G14a3EgKyNeeiY"
}
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
     ip ="${profitbricks_ipblock.webserver_ip.ips[0]}"
  }
}`
