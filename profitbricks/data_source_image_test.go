package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func TestAccDataSourceImage_basic(t *testing.T) {
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
					testAccCheckProfitBricksImageAttributes("data.profitbricks_image.img", "Ubuntu-16.04"),
					resource.TestCheckResourceAttr("data.profitbricks_image.img", "type", "HDD"),
				),
			},
		},
	})

}

func testAccCheckProfitBricksImageAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksImageAttributes: Not found: %s", n)
		}

		if !strings.Contains(rs.Primary.Attributes["name"], name) {
			return fmt.Errorf("Bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

const testAccDataSourceProfitBricksImage_basic = `
	data "profitbricks_image" "img" {
	  name = "Ubuntu"
	  type = "HDD"
	  version = "16"
	  location = "us/las"
	}
`
