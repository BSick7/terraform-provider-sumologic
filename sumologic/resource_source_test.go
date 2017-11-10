package sumologic

import (
	"fmt"
	"strconv"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckSourceExists(n string, res *api.Source) resource.TestCheckFunc {
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
		cid, _ := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		source, err := client.Collectors().Sources(cid).Get(id)
		if err != nil {
			return err
		}

		if source == nil {
			return fmt.Errorf("source not found")
		}

		*res = *source

		return nil
	}
}

func testAccCheckSourceDestroy(resourceType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			client := testAccProvider.Meta().(*api.Client)
			id, _ := strconv.Atoi(rs.Primary.ID)
			cid, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
			if err != nil {
				return fmt.Errorf("invalid collector_id: %s", rs.Primary.Attributes["collector_id"])
			}

			if exists, err := client.Collectors().Sources(cid).Exists(id); err != nil {
				return fmt.Errorf("error checking for source destruction: %s", err)
			} else if exists {
				return fmt.Errorf("source was not destroyed")
			}
		}

		return nil
	}
}
