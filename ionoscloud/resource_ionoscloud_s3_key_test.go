package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudS3Key_Basic(t *testing.T) {
	var s3Key ionoscloud.S3Key
	s3KeyName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosClouds3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosClouds3KeyConfigBasic, s3KeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosClouds3KeyExists("ionoscloud_s3_key.example", &s3Key),
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example", "secret_key"),
					resource.TestCheckResourceAttr("ionoscloud_s3_key.example", "active", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckIonosClouds3KeyConfigUpdate, s3KeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosClouds3KeyExists("ionoscloud_s3_key.example", &s3Key),
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example", "secret_key"),
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example", "active"),
				),
			},
		},
	})
}

func testAccCheckDIonosClouds3KeyDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_key" {
			continue
		}

		_, err := client.GetS3Key(rs.Primary.Attributes["user_id"], rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("S3 Key still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch S3 Key %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosClouds3KeyExists(n string, s3Key *ionoscloud.S3Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
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

const testAccCheckIonosClouds3KeyConfigBasic = `
resource "ionoscloud_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "terraform-s3-acc-tester2@profitbricks.com"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_s3_key" "%s" {
  user_id    = ionoscloud_user.example.id
  active     = true
}`

const testAccCheckIonosClouds3KeyConfigUpdate = `
resource "ionoscloud_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "terraform-s3-acc-tester2@profitbricks.com"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_s3_key" "%s" {
  user_id    = ionoscloud_user.example.id
  active     = false
}`
