package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
	"math/rand"
	"strconv"
	"time"
)

func TestAccProfitBricksUser_Basic(t *testing.T) {
	var user profitbricks.User
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	email := "terraform_test" + strconv.Itoa(r1.Intn(1000)) + "@go.com"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksUserConfig_basic, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksUserExists("profitbricks_user.user", &user),
					testAccCheckProfitBricksUserAttributes("profitbricks_user.user", "terraform"),
					resource.TestCheckResourceAttr("profitbricks_user.user", "first_name", "terraform"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckProfitbricksUserConfig_update, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksUserAttributes("profitbricks_user.user", "updated"),
					resource.TestCheckResourceAttr("profitbricks_user.user", "first_name", "updated"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksUserDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		resp := profitbricks.GetUser(rs.Primary.ID)

		if resp.StatusCode < 299 {
			return fmt.Errorf("user still exists %s %s", rs.Primary.ID, resp.Response)
		}
	}

	return nil
}

func testAccCheckProfitBricksUserAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksUserAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["first_name"] != name {
			return fmt.Errorf("Bad first_name: %s", rs.Primary.Attributes["first_name"])
		}

		return nil
	}
}

func testAccCheckProfitBricksUserExists(n string, user *profitbricks.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckProfitBricksUserExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		founduser := profitbricks.GetUser(rs.Primary.ID)

		if founduser.StatusCode != 200 {
			return fmt.Errorf("Error occured while fetching User: %s", rs.Primary.ID)
		}
		if founduser.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		user = &founduser

		return nil
	}
}

const testAccCheckProfitbricksUserConfig_basic = `

resource "profitbricks_group" "group" {
  name = "terraform user group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "profitbricks_user" "user" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}`

const testAccCheckProfitbricksUserConfig_update = `


resource "profitbricks_user" "user" {
  first_name = "updated"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "profitbricks_group" "group" {
  name = "terraform user group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  user_id="${profitbricks_user.user.id}"
}
`
