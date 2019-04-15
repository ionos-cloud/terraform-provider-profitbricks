package profitbricks

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceImage_basic(t *testing.T) {
	r, _ := regexp.Compile("Ubuntu-16.04")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceProfitBricksImage_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_image.img", "location", "us/las"),
					resource.TestMatchResourceAttr("data.profitbricks_image.img", "name", r),
					resource.TestCheckResourceAttr("data.profitbricks_image.img", "type", "HDD"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksImage_basic = `
	data "profitbricks_image" "img" {
	  name = "Ubuntu"
	  type = "HDD"
	  version = "16"
	  location = "us/las"
	}
`
