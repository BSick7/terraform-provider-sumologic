package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHTTPSource_Basic(t *testing.T) {
	var source api.Source

	testCheck := func(s *terraform.State) error {
		if source.ID <= 0 {
			return fmt.Errorf("expected source to be created")
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceDestroy("sumologic_http_source"),
		Steps: []resource.TestStep{
			{
				Config: testAccHTTPSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists("sumologic_http_source.source1", &source),
					testCheck,
				),
			},
		},
	})
}

var testAccHTTPSourceBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector2" {
  name        = "tf-acc-collector2"
  description = "Collector 2 (TF Acceptance Test)"
}
resource "sumologic_http_source" "source1" {
  collector_id = "${sumologic_hosted_collector.collector2.id}"
  name         = "tf-acc-source1"
  description  = "Source 1 (TF Acceptance Test)"
}
`)
