package ionoscloud

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
				Config: testAccDataSourceIonosCloudServerCreateResources,
			},
			{
				Config: testAccDataSourceIonosCloudServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_server.test_server", "name", "test_ds_server"),
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
				Config: testAccDataSourceIonosCloudServerCreateResources,
			},
			{
				Config: testAccDataSourceIonosCloudServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_server.test_server", "name", "test_ds_server"),
				),
			},
		},
	})

}

const testAccDataSourceIonosCloudServerCreateResources = `
resource "ionoscloud_datacenter" "test_ds_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}
resource "ionoscloud_server"" "test_ds_server" {
  depends_on        = [ionoscloud_datacenter.test_ds_server]
  datacenter_id     = ionoscloud_datacenter.test_ds_server.id
  name              = "test_ds_server"
}
`

const testAccDataSourceIonosCloudServerMatchId = `
resource "ionoscloud_datacenter" "test_ds_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}

resource "ionoscloud_server" "test_ds_server" {
  depends_on        = [ionoscloud_datacenter.test_ds_server]
  datacenter_id     = ionoscloud_datacenter.test_ds_server.id
  name              = "test_ds_server"
}

data "ionoscloud_server" "test_server" {
  datacenter_id = ionoscloud_datacenter.test_ds_server.id
  id			= ionoscloud_server.test_ds_server.id
}
`

const testAccDataSourceIonosCloudServerMatchName = `
resource "ionoscloud_datacenter" "test_ds_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}

resource "ionoscloud_server" "test_ds_server" {
  depends_on        = [ionoscloud_datacenter.test_ds_server]
  datacenter_id     = ionoscloud_datacenter.test_ds_server.id
  name              = "test_ds_server"
  public            = true
}

data "ionoscloud_server" "test_server" {
  datacenter_id = ionoscloud_datacenter.test_ds_server.id
  name			= "test_ds_"
}
`
