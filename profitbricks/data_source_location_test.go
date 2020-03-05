package profitbricks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLocation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceProfitBricksLocation_basic,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("data.profitbricks_location.loc", "id", "de/fkb"),
					resource.TestCheckResourceAttr("data.profitbricks_location.loc", "name", "karlsruhe"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksLocation_basic = `
	data "profitbricks_location" "loc" {
	  name = "karlsruhe"
	  feature = "SSD"
	}
	`
