package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudWatchSource_Basic(t *testing.T) {
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
		CheckDestroy: testAccCheckSourceDestroy("sumologic_cloudwatch_source"),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudWatchSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists("sumologic_cloudwatch_source.source1", &source),
					testCheck,
				),
			},
		},
	})
}

var testAccCloudWatchSourceBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector1" {
  name        = "tf-acc-collector1-http"
  description = "Collector HTTP 1 (TF Acceptance Test)"
}
resource "sumologic_cloudwatch_source" "source1" {
  collector_id = "${sumologic_hosted_collector.collector1.id}"
  name         = "tf-acc-cloudwatch-source1"
  description  = "CloudWatch Source 1 (TF Acceptance Test)"
  aws_access_key = "stub-access-key"
  aws_secret_key = "stub-secret-key"
  aws_bucket     = "stub-bucket"
  scan_interval  = "1m0s"
}
`)
