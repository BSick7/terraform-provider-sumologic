package sumologic

import (
	"fmt"
	"strconv"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func importCollector(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// If identifier is an int, we will import by collector id
	// If identifier is a string, we will import by name
	id, err := findCollectorIDByIdentifier(meta.(*api.Client), d.Id())
	d.SetId(fmt.Sprintf("%d", id))
	return []*schema.ResourceData{d}, err
}

func findCollectorIDByIdentifier(client *api.Client, identifier string) (int, error) {
	id, err := strconv.Atoi(identifier)
	if err == nil {
		return id, nil
	}

	limit := 100
	for i := 0; i <= 10000; i += limit {
		collectors, err := client.Collectors().List(i, limit)
		if err != nil {
			return 0, err
		}
		for _, collector := range collectors {
			if collector.Name == identifier {
				return collector.ID, nil
			}
		}
	}
	return 0, nil
}
