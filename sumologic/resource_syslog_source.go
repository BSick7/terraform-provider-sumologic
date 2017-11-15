package sumologic

import (
	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSyslogSource() *schema.Resource {
	sch := resourceSourceSchema()

	return &schema.Resource{
		Create: resourceSyslogSourceCreate,
		Read:   resourceSyslogSourceRead,
		Update: resourceSyslogSourceUpdate,
		Delete: resourceSourceDelete,
		Exists: resourceSourceExists,
		Importer: &schema.ResourceImporter{
			State: importSource,
		},
		Schema: sch,
	}
}

func resourceSyslogSourceCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceSourceCreate(d, meta, func(source *api.SourceCreate) error {
		source.SourceType = "Cloudsyslog"
		return nil
	})
	if err != nil {
		return err
	}
	if err := resourceHTTPSourceUpdate(d, meta); err != nil {
		return err
	}
	return resourceHTTPSourceRead(d, meta)
}

func resourceSyslogSourceRead(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceRead(d, meta, func(source *api.Source) error {
		return nil
	})
}

func resourceSyslogSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceUpdate(d, meta, func(source *api.Source) error {
		source.SourceType = "Cloudsyslog"
		return nil
	})
}
