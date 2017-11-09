package sumologic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func importSource(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// terraform import <resource> <collector-id/name>:<source-id/name>
	// The id specified will include collector id/name and source id/name
	// They are delimited by ":"
	tokens := strings.Split(d.Id(), ":")
	cid, err := findCollectorIDByIdentifier(meta.(*api.Client), tokens[0])
	if err != nil {
		return nil, err
	}
	sid, err := findSourceByName(meta.(*api.Client), cid, tokens[1])
	d.Set("collector_id", cid)
	d.SetId(fmt.Sprintf("%d", sid))
	return []*schema.ResourceData{d}, err
}

func findSourceByName(client *api.Client, collectorID int, identifier string) (int, error) {
	id, err := strconv.Atoi(identifier)
	if err == nil {
		return id, nil
	}

	sources, err := client.Collectors().Sources(collectorID).List()
	if err != nil {
		return 0, err
	}
	for _, source := range sources {
		if source.Name == identifier {
			return source.ID, nil
		}
	}
	return 0, nil
}
