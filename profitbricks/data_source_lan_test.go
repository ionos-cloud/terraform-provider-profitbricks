package profitbricks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLan_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksLanCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksLanMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_lan.test_lan", "name", "test_ds_lan"),
				),
			},
		},
	})
}

func TestAccDataSourceLan_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksLanCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksLanMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_lan.test_lan", "name", "test_ds_lan"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksLanCreateResources = `
resource "profitbricks_datacenter" "test_ds_lan" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}
resource "profitbricks_lan" "test_ds_lan" {
  depends_on        = [profitbricks_datacenter.test_ds_lan]
  datacenter_id     = profitbricks_datacenter.test_ds_lan.id
  name              = "test_ds_lan"
  public            = true
}
`

const testAccDataSourceProfitBricksLanMatchId = `
resource "profitbricks_datacenter" "test_ds_lan" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}

resource "profitbricks_lan" "test_ds_lan" {
  depends_on        = [profitbricks_datacenter.test_ds_lan]
  datacenter_id     = profitbricks_datacenter.test_ds_lan.id
  name              = "test_ds_lan"
  public            = true
}

data "profitbricks_lan" "test_lan" {
  datacenter_id = profitbricks_datacenter.test_ds_lan.id
  id			= profitbricks_lan.test_ds_lan.id
}
`

const testAccDataSourceProfitBricksLanMatchName = `
resource "profitbricks_datacenter" "test_ds_lan" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}

resource "profitbricks_lan" "test_ds_lan" {
  depends_on        = [profitbricks_datacenter.test_ds_lan]
  datacenter_id     = profitbricks_datacenter.test_ds_lan.id
  name              = "test_ds_lan"
  public            = true
}

data "profitbricks_lan" "test_lan" {
  datacenter_id = profitbricks_datacenter.test_ds_lan.id
  name			= "test_ds_"
}
`
