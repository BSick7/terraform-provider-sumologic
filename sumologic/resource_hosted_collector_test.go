package sumologic

import (
	"fmt"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHostedCollector_Basic(t *testing.T) {
	var collector api.Collector

	testCheck := func(s *terraform.State) error {
		if collector.ID <= 0 {
			return fmt.Errorf("expected collector to be created")
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy("sumologic_hosted_collector"),
		Steps: []resource.TestStep{
			{
				Config: testAccHostedCollectorBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_hosted_collector.collector1", &collector),
					testCheck,
				),
			},
		},
	})
}

var testAccHostedCollectorBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector1" {
  name        = "tf-acc-collector1"
  description = "Collector 1 (TF Acceptance Test)"
}
`)
