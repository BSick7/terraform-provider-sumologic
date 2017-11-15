package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSyslogSource_Basic(t *testing.T) {
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
		CheckDestroy: testAccCheckSourceDestroy("sumologic_syslog_source.source1"),
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists("sumologic_syslog_source.source1", &source),
					testCheck,
				),
			},
		},
	})
}

var testAccSyslogSourceBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector1" {
  name        = "tf-acc-collector1-syslog"
  description = "Collector Syslog 1 (TF Acceptance Test)"
}
resource "sumologic_syslog_source" "source1" {
  collector_id   = "${sumologic_hosted_collector.collector1.id}"
  name           = "tf-acc-syslog-source1"
  description    = "Syslog Source 1 (TF Acceptance Test)"
}
`)
