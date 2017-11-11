package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLocalFileSource_Basic(t *testing.T) {
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
		CheckDestroy: testAccCheckSourceDestroy("sumologic_local_file_source"),
		Steps: []resource.TestStep{
			{
				Config: testAccLocalFileSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists("sumologic_local_file_source.source1", &source),
					testCheck,
				),
			},
		},
	})
}

var testAccLocalFileSourceBasicConfig = fmt.Sprint(`
resource "sumologic_installed_collector" "collector1" {
  name        = "tf-acc-collector1-lf"
  description = "Collector LF 1 (TF Acceptance Test)"
}
resource "sumologic_local_file_source" "source1" {
  collector_id = "${sumologic_hosted_collector.collector1.id}"
  name            = "tf-acc-lf-source1"
  description     = "Local File Source 1 (TF Acceptance Test)"
  path_expression = "/var/log/syslog"
}
`)
