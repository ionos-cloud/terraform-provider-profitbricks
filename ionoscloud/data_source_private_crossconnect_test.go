package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcePcc_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIonosCloudPccCreateResources,
			},
			{
				Config: testAccDataSourceIonosCloudPccMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_private_crossconnect.test_pcc", "name", "test_ds_pcc"),
					resource.TestCheckResourceAttr("data.ionoscloud_private_crossconnect.test_pcc", "description", "test_ds_pcc description"),
				),
			},
		},
	})
}

func TestAccDataSourcePcc_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIonosCloudPccCreateResources,
			},
			{
				Config: testAccDataSourceIonosCloudPccMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_private_crossconnect.test_pcc", "name", "test_ds_pcc"),
					resource.TestCheckResourceAttr("data.ionoscloud_private_crossconnect.test_pcc", "description", "test_ds_pcc description"),
				),
			},
		},
	})

}

const testAccDataSourceIonosCloudPccCreateResources = `
resource "ionoscloud_private_crossconnect" "test_ds_pcc" {
  name              = "test_ds_pcc"
  description		= "test_ds_pcc description"
}
`

const testAccDataSourceIonosCloudPccMatchId = `
resource "ionoscloud_private_crossconnect" "test_ds_pcc" {
  name              = "test_ds_pcc"
  description		= "test_ds_pcc description"
}

data "ionoscloud_private_crossconnect" "test_pcc" {
  id			= ionoscloud_private_crossconnect.test_ds_pcc.id
}
`

const testAccDataSourceIonosCloudPccMatchName = `
resource "ionoscloud_private_crossconnect" "test_ds_pcc" {
  name              = "test_ds_pcc"
  description		= "test_ds_pcc description"
}

data "ionoscloud_private_crossconnect" "test_pcc" {
  name			= "test_ds_"
}
`
