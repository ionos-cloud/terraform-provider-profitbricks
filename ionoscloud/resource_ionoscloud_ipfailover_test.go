package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudLanIPFailover_Basic(t *testing.T) {
	var lan ionoscloud.Lan
	var ipfailover ionoscloud.IPFailover

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
		CheckDestroy: testAccCheckDIonosCloudLanIPFailoverDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksLanIPFailoverConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanIPFailoverGroupExists("ionoscloud_ipfailover.failovertest", &lan, &ipfailover),
				),
			},
			{
				Config: testAccCheckProfitbricksLanIPFailoverConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("ionoscloud_ipfailover.failovertest"),
				),
			},
		},
	})
}

func testAccCheckLanIPFailoverGroupExists(n string, lan *ionoscloud.Lan, failover *ionoscloud.IPFailover) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID` is set")
		}

		lanId := rs.Primary.Attributes["lan_id"]
		nicUuid := rs.Primary.Attributes["nicuuid"]

		lan, err := client.GetLan(rs.Primary.Attributes["datacenter_id"], lanId)
		if err != nil {
			return fmt.Errorf("Lan %s not found.", lanId)
		}

		if lan.Properties.IPFailover == nil {
			return fmt.Errorf("Lan %s has no failover groups.", lanId)
		}
		found := false
		for _, fo := range *lan.Properties.IPFailover {
			if fo.NicUUID == nicUuid {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Expected NIC %s to be a part of a failover group, but not found in lans %s failover groups", nicUuid, lanId)
		}

		return nil
	}
}

func testAccCheckDIonosCloudLanIPFailoverDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_ipfailover" {
			continue
		}
		nicUuid := rs.Primary.Attributes["nicuuid"]
		lan, err := client.GetLan(rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["lan_id"])

		if err != nil {
			return fmt.Errorf("An error occured while fetching a Lan ID %s %s", rs.Primary.Attributes["lan_id"], err)
		}

		found := false
		for _, fo := range *lan.Properties.IPFailover {
			if fo.NicUUID == nicUuid {
				found = true
			}
		}
		if found {
			_, err := client.DeleteDatacenter(rs.Primary.Attributes["datacenter_id"])
			if err != nil {
				return fmt.Errorf("IP failover group with nicId %s still exists %s %s, removing datacenter....", nicUuid, rs.Primary.ID, err)
			}
		}
	}

	return nil
}

const testAccCheckProfitbricksLanIPFailoverConfig_basic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 1
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "centos:latest"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan1.id}"
    dhcp = true
    firewall_active = true
     ip ="${ionoscloud_ipblock.webserver_ip.ips[0]}"
  }
}
resource "ionoscloud_ipfailover" "failovertest" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  lan_id="${ionoscloud_lan.webserver_lan1.id}"
  ip ="${ionoscloud_ipblock.webserver_ip.ips[0]}"
  nicuuid= "${ionoscloud_server.webserver.primary_nic}"
}`

const testAccCheckProfitbricksLanIPFailoverConfig_update = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 1
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name = "centos:latest"
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
     ip ="${ionoscloud_ipblock.webserver_ip.ips[0]}"
  }
}`
