package profitbricks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceServer_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksServerCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_server.test_server", "name", "test_datasource_server"),
				),
			},
		},
	})
}

func TestAccDataSourceServer_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksServerCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_server.test_server", "name", "test_datasource_server"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksServerCreateResources = `
resource "profitbricks_datacenter" "test_datasource_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}
resource "profitbricks_server" "test_datasource_server" {
  depends_on        = [profitbricks_datacenter.test_datasource_server]
  datacenter_id     = profitbricks_datacenter.test_datasource_server.id
  name              = "test_datasource_server"
  cores             = 2
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_XEON"

  image_name        = "ubuntu:20.04"
  image_password    = "foobar123456"

  volume {
    size            =   "40"
    disk_type       =   "HDD"
  }

  nic {
    lan             = 1
  }
}
`

const testAccDataSourceProfitBricksServerMatchId = `
resource "profitbricks_datacenter" "test_datasource_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}

resource "profitbricks_server" "test_datasource_server" {
  depends_on        = [profitbricks_datacenter.test_datasource_server]
  datacenter_id     = profitbricks_datacenter.test_datasource_server.id
  name              = "test_datasource_server"
  cores             = 2
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_XEON"

  image_name        = "ubuntu:20.04"
  image_password    = "foobar123456"

  volume {
    size            =   "40"
    disk_type       =   "HDD"
  }

  nic {
    lan             = 1
  }
}

data "profitbricks_server" "test_server" {
  datacenter_id = profitbricks_datacenter.test_datasource_server.id
  id			= profitbricks_server.test_datasource_server.id
}
`

const testAccDataSourceProfitBricksServerMatchName = `
resource "profitbricks_datacenter" "test_datasource_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}

resource "profitbricks_server" "test_datasource_server" {
  depends_on        = [profitbricks_datacenter.test_datasource_server]
  datacenter_id     = profitbricks_datacenter.test_datasource_server.id
  name              = "test_datasource_server"

  cores             = 2
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_XEON"

  image_name        = "ubuntu:20.04"
  image_password    = "foobar123456"

  volume {
    size            =   "40"
    disk_type       =   "HDD"
  }

  nic {
    lan             = 1
  }
}

data "profitbricks_server" "test_server" {
  datacenter_id = profitbricks_datacenter.test_datasource_server.id
  name			= "test_datasource_"
}
`
