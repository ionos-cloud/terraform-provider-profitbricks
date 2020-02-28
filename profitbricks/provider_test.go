package profitbricks

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"profitbricks": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	pbUsername := os.Getenv("PROFITBRICKS_USERNAME")
	pbPassword := os.Getenv("PROFITBRICKS_PASSWORD")
	pbToken := os.Getenv("PROFITBRICKS_TOKEN")
	if pbToken == "" {
		if pbUsername == "" || pbPassword == "" {
			t.Fatal("PROFITBRICKS_USERNAME/PROFITBRICKS_PASSWORD or PROFITBRICKS_TOKEN must be set for acceptance tests")
		}
	} else {
		if pbUsername != "" || pbPassword != "" {
			t.Fatal("PROFITBRICKS_USERNAME/PROFITBRICKS_PASSWORD or PROFITBRICKS_TOKEN must be set for acceptance tests")
		}

	}
}
