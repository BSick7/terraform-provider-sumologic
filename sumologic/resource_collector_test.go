package sumologic

import (
	"fmt"
	"strconv"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckCollectorExists(n string, res *api.Collector) resource.TestCheckFunc {
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
			return fmt.Errorf("collector not found")
		}

		*res = *collector

		return nil
	}
}

func testAccCheckCollectorDestroy(resourceType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			client := testAccProvider.Meta().(*api.Client)
			id, _ := strconv.Atoi(rs.Primary.ID)
			if exists, err := client.Collectors().Exists(id); err != nil {
				return fmt.Errorf("error checking for collector destruction: %s", err)
			} else if exists {
				return fmt.Errorf("collector was not destroyed")
			}
		}

		return nil
	}
}
