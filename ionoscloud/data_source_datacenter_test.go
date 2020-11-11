package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDatacenter_matching(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceIonosCloudDataCenter_matching,
			},
			{
				Config: testAccDataSourceIonosCloudDataCenter_matchingWithDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_datacenter.foobar", "name", "test_name"),
					resource.TestCheckResourceAttr("data.ionoscloud_datacenter.foobar", "location", "us/las"),
				),
			},
		},
	})

}

const testAccDataSourceIonosCloudDataCenter_matching = `
resource "ionoscloud_datacenter" "foobar" {
    name       = "test_name"
    location = "us/las"
}
`

const testAccDataSourceIonosCloudDataCenter_matchingWithDataSource = `
resource "ionoscloud_datacenter" "foobar" {
    name       = "test_name"
    location = "us/las"
}

data "ionoscloud_datacenter" "foobar" {
    name = "${ionoscloud_datacenter.foobar.name}"
    location = "us/las"
}`
