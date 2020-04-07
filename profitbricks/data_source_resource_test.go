package profitbricks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksResource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_resource.res", "resource_type", "datacenter"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksResource_basic = `
resource "profitbricks_datacenter" "foobar" {
  name       = "test_name"
  location = "us/las"
}

data "profitbricks_resource" "res" {
  resource_type = "datacenter"
  resource_id="${profitbricks_datacenter.foobar.id}"
}`
