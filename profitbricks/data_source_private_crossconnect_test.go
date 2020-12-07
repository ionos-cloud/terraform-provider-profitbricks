package profitbricks

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
				Config: testAccDataSourceProfitBricksPccCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksPccMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_private_crossconnect.test_pcc", "name", "test_ds_pcc"),
					resource.TestCheckResourceAttr("data.profitbricks_private_crossconnect.test_pcc", "description", "test_ds_pcc description"),
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
				Config: testAccDataSourceProfitBricksPccCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksPccMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_private_crossconnect.test_pcc", "name", "test_ds_pcc"),
					resource.TestCheckResourceAttr("data.profitbricks_private_crossconnect.test_pcc", "description", "test_ds_pcc description"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksPccCreateResources = `
resource "profitbricks_private_crossconnect" "test_ds_pcc" {
  name              = "test_ds_pcc"
  description		= "test_ds_pcc description"
}
`

const testAccDataSourceProfitBricksPccMatchId = `
resource "profitbricks_private_crossconnect" "test_ds_pcc" {
  name              = "test_ds_pcc"
  description		= "test_ds_pcc description"
}

data "profitbricks_private_crossconnect" "test_pcc" {
  id			= profitbricks_private_crossconnect.test_ds_pcc.id
}
`

const testAccDataSourceProfitBricksPccMatchName = `
resource "profitbricks_private_crossconnect" "test_ds_pcc" {
  name              = "test_ds_pcc"
  description		= "test_ds_pcc description"
}

data "profitbricks_private_crossconnect" "test_pcc" {
  name			= "test_ds_"
}
`
