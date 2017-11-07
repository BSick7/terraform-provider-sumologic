package sumologic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sumologic": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SUMO_ACCESS_ID"); v == "" {
		t.Fatal("SUMO_ACCESS_ID must be set for acceptance tests")
	}
	if v := os.Getenv("SUMO_ACCESS_KEY"); v == "" {
		t.Fatal("SUMO_ACCESS_KEY must be set for acceptance tests")
	}
}
