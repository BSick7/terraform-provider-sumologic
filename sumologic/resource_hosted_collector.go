package sumologic

import (
	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHostedCollector() *schema.Resource {
	sch := resourceCollectorSchema()

	return &schema.Resource{
		Create: resourceHostedCollectorCreate,
		Read:   resourceHostedCollectorRead,
		Update: resourceHostedCollectorUpdate,
		Delete: resourceCollectorDelete,
		Exists: resourceCollectorExists,
		Importer: &schema.ResourceImporter{
			State: importCollector,
		},
		Schema: sch,
	}
}

func resourceHostedCollectorCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceCollectorCreate(d, meta, func(collector *api.CollectorCreate) error {
		collector.CollectorType = "Hosted"
		return nil
	})
	if err != nil {
		return err
	}
	if err := resourceHostedCollectorUpdate(d, meta); err != nil {
		return err
	}

	return resourceHostedCollectorRead(d, meta)
}

func resourceHostedCollectorRead(d *schema.ResourceData, meta interface{}) error {
	return resourceCollectorRead(d, meta, func(collector *api.Collector) error {
		return nil
	})
}

func resourceHostedCollectorUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCollectorUpdate(d, meta, func(collector *api.Collector) error {
		collector.CollectorType = "Hosted"
		return nil
	})
}
