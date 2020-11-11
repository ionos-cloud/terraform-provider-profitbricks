package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccIonosCloudbackupUnit_Basic(t *testing.T) {
	var backupUnit ionoscloud.BackupUnit
	backupUnitName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDIonosCloudbackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIonosCloudbackupUnitConfigBasic, backupUnitName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudbackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "name", backupUnitName),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example@ionos.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword123$"),
				),
			},
			{
				Config: testAccCheckIonosCloudbackupUnitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIonosCloudbackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example-updated@ionos.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword1234$"),
				),
			},
		},
	})
}

func testAccCheckDIonosCloudbackupUnitDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_backup_unit" {
			continue
		}

		_, err := client.GetBackupUnit(rs.Primary.ID)

		if apiError, ok := err.(ionoscloud.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("backup unit still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch backup unit %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIonosCloudbackupUnitExists(n string, backupUnit *ionoscloud.BackupUnit) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundBackupUnit, err := client.GetBackupUnit(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching backup unit: %s", rs.Primary.ID)
		}
		if foundBackupUnit.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		backupUnit = foundBackupUnit

		return nil
	}
}

const testAccCheckIonosCloudbackupUnitConfigBasic = `
resource "ionoscloud_backup_unit" "example" {
  name        = "%s"
	password    = "DemoPassword123$"
  email       = "example@ionos.com"
}`

const testAccCheckIonosCloudbackupUnitConfigUpdate = `
resource "ionoscloud_backup_unit" "example" {
	name        = "example"
	email       = "example-updated@ionos.com"
	password    = "DemoPassword1234$"
}`
