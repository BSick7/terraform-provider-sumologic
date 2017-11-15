package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBucketSource_All(t *testing.T) {
	var source api.Source

	testCheck := func(s *terraform.State) error {
		if source.ID <= 0 {
			return fmt.Errorf("expected source to be created")
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckSourcesDestroy([]string{
			"sumologic_cloudfront_source.source1",
			"sumologic_cloudtrail_source.source1",
			"sumologic_elb_source.source1",
			"sumologic_s3_source.source1",
			"sumologic_s3_audit_source.source1",
		}),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists("sumologic_cloudfront_source.source1", &source),
					testCheck,
					testAccCheckSourceExists("sumologic_cloudtrail_source.source1", &source),
					testCheck,
					testAccCheckSourceExists("sumologic_elb_source.source1", &source),
					testCheck,
					testAccCheckSourceExists("sumologic_s3_source.source1", &source),
					testCheck,
					testAccCheckSourceExists("sumologic_s3_audit_source.source1", &source),
					testCheck,
				),
			},
		},
	})
}

var testAccBucketSourceBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector1" {
  name        = "tf-acc-collector1-bucket"
  description = "Collector Bucket 1 (TF Acceptance Test)"
}
resource "sumologic_cloudfront_source" "source1" {
  collector_id   = "${sumologic_hosted_collector.collector1.id}"
  name           = "tf-acc-cloudfront-source1"
  description    = "CloudFront Source 1 (TF Acceptance Test)"
  aws_access_key = "stub-access-key"
  aws_secret_key = "stub-secret-key"
  aws_bucket     = "stub-bucket"
  scan_interval  = "1m0s"
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
resource "sumologic_elb_source" "source1" {
  collector_id   = "${sumologic_hosted_collector.collector1.id}"
  name           = "tf-acc-elb-source1"
  description    = "ELB Source 1 (TF Acceptance Test)"
  aws_access_key = "stub-access-key"
  aws_secret_key = "stub-secret-key"
  aws_bucket     = "stub-bucket"
  scan_interval  = "1m0s"
}
resource "sumologic_s3_source" "source1" {
  collector_id   = "${sumologic_hosted_collector.collector1.id}"
  name           = "tf-acc-s3-source1"
  description    = "S3 Source 1 (TF Acceptance Test)"
  aws_access_key = "stub-access-key"
  aws_secret_key = "stub-secret-key"
  aws_bucket     = "stub-bucket"
  scan_interval  = "1m0s"
}
resource "sumologic_s3_audit_source" "source1" {
  collector_id   = "${sumologic_hosted_collector.collector1.id}"
  name           = "tf-acc-s3-audit-source1"
  description    = "S3 Audit Source 1 (TF Acceptance Test)"
  aws_access_key = "stub-access-key"
  aws_secret_key = "stub-secret-key"
  aws_bucket     = "stub-bucket"
  scan_interval  = "1m0s"
}
`)
