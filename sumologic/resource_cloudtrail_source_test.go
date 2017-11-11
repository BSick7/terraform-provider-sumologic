package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudTrailSource_Basic(t *testing.T) {
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
		CheckDestroy: testAccCheckSourceDestroy("sumologic_cloudtrail_source"),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTrailSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists("sumologic_cloudtrail_source.source1", &source),
					testCheck,
				),
			},
		},
	})
}

var testAccCloudTrailSourceBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector1" {
  name        = "tf-acc-collector1-cloudtrail"
  description = "Collector CloudTrail 1 (TF Acceptance Test)"
}
resource "sumologic_cloudtrail_source" "source1" {
  collector_id   = "${sumologic_hosted_collector.collector1.id}"
  name           = "tf-acc-cloudtrail-source1"
  description    = "CloudTrail Source 1 (TF Acceptance Test)"
  aws_access_key = "stub-access-key"
  aws_secret_key = "stub-secret-key"
  aws_bucket     = "stub-bucket"
  scan_interval  = "1m0s"
}
`)
