package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccProfitBricksDataCenter_Basic(t *testing.T) {
	var datacenter profitbricks.Datacenter
	dc_name := "datacenter-test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksDatacenterConfig_basic, dc_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksDatacenterExists("profitbricks_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("profitbricks_datacenter.foobar", "name", dc_name),
				),
			},
			{
				Config: testAccCheckProfitBricksDatacenterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksDatacenterExists("profitbricks_datacenter.foobar", &datacenter),
					resource.TestCheckResourceAttr("profitbricks_datacenter.foobar", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksDatacenterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_datacenter" {
			continue
		}

		_, err := client.GetDatacenter(rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("DataCenter still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetching DataCenter %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckProfitBricksDatacenterExists(n string, datacenter *profitbricks.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundDC, err := client.GetDatacenter(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching DC: %s", rs.Primary.ID)
		}
		if foundDC.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		datacenter = foundDC

		return nil
	}
}

const testAccCheckProfitBricksDatacenterConfig_basic = `
resource "profitbricks_datacenter" "foobar" {
	name       = "%s"
	location = "us/las"
}`

const testAccCheckProfitBricksDatacenterConfig_update = `
resource "profitbricks_datacenter" "foobar" {
	name       =  "updated"
	location = "us/las"
}`
