package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccProfitBricksS3Key_ImportBasic(t *testing.T) {
	resourceName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBrickss3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBrickss3KeyImportConfigBasic, resourceName),
			},
			{
				ResourceName:            fmt.Sprintf("profitbricks_s3_key.%s", resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccProfitBricksS3KeyImportStateID,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccProfitBricksS3KeyImportStateID(s *terraform.State) (string, error) {
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_s3_key" {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["user_id"], rs.Primary.ID)
	}

	return importID, nil
}

const testAccCheckProfitBrickss3KeyImportConfigBasic = `

resource "profitbricks_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "terraform-s3-import-acc-tester2@ionos.com"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "profitbricks_s3_key" "%s" {
  user_id    = profitbricks_user.example.id
  active     = true
}`
