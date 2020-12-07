package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccProfitBricksS3Key_Basic(t *testing.T) {
	var s3Key profitbricks.S3Key
	s3KeyName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBrickss3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBrickss3KeyConfigBasic, s3KeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBrickss3KeyExists("profitbricks_s3_key.example", &s3Key),
					resource.TestCheckResourceAttrSet("profitbricks_s3_key.example", "secret_key"),
					resource.TestCheckResourceAttr("profitbricks_s3_key.example", "active", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckProfitBrickss3KeyConfigUpdate, s3KeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBrickss3KeyExists("profitbricks_s3_key.example", &s3Key),
					resource.TestCheckResourceAttrSet("profitbricks_s3_key.example", "secret_key"),
					resource.TestCheckResourceAttrSet("profitbricks_s3_key.example", "active"),
				),
			},
		},
	})
}

func testAccCheckDProfitBrickss3KeyDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_s3_key" {
			continue
		}

		_, err := client.GetS3Key(rs.Primary.Attributes["user_id"], rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("S3 Key still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch S3 Key %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckProfitBrickss3KeyExists(n string, s3Key *profitbricks.S3Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundS3Key, err := client.GetS3Key(rs.Primary.Attributes["user_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching S3 Key: %s", rs.Primary.ID)
		}

		if foundS3Key.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		s3Key = foundS3Key

		return nil
	}
}

const testAccCheckProfitBrickss3KeyConfigBasic = `
resource "profitbricks_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "terraform-s3-acc-tester2@profitbricks.com"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "profitbricks_s3_key" "%s" {
  user_id    = profitbricks_user.example.id
  active     = true
}`

const testAccCheckProfitBrickss3KeyConfigUpdate = `
resource "profitbricks_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "terraform-s3-acc-tester2@profitbricks.com"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "profitbricks_s3_key" "%s" {
  user_id    = profitbricks_user.example.id
  active     = false
}`
