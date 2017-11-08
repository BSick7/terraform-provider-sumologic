package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCollector_Basic(t *testing.T) {
	var collector *api.Collector

	testCheck := func(s *terraform.State) error {
		if collector == nil {
			return fmt.Errorf("expected collector to be created")
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSLHostedCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumoLogicCollectorBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumologicCollectorExists("sumologic_hosted_collector.collector1", collector),
					testCheck,
				),
			},
		},
	})
}

func testAccCheckSumologicCollectorExists(n string, res *api.Collector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		id, _ := strconv.Atoi(rs.Primary.ID)
		collector, err := client.Collectors().Get(id)
		if err != nil {
			return err
		}

		if collector == nil {
			return fmt.Errorf("subnet not found")
		}

		*res = *collector

		return nil
	}
}

func testAccCheckSLHostedCollectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_hosted_collector" {
			continue
		}

		client := testAccProvider.Meta().(*api.Client)
		id, _ := strconv.Atoi(rs.Primary.ID)
		collector, err := client.Collectors().Get(id)
		if err != nil {
			return err
		}

		if collector != nil {
			return fmt.Errorf("collector was not destroyed. %+v", collector)
		}
	}

	return nil
}

var testAccSumoLogicCollectorBasicConfig = fmt.Sprint(`
resource "sumologic_hosted_collector" "collector1" {
  name        = "collector1-testacc"
  description = "Collector 1 (Acceptance Test)"
}
`)
