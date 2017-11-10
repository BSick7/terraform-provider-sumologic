package sumologic

import (
	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHTTPSource() *schema.Resource {
	sch := resourceSourceSchema()
	sch["message_per_request"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Enable One Message Per Request",
	}

	return &schema.Resource{
		Create: resourceHTTPSourceCreate,
		Read:   resourceHTTPSourceRead,
		Update: resourceHTTPSourceUpdate,
		Delete: resourceSourceDelete,
		Exists: resourceSourceExists,
		Importer: &schema.ResourceImporter{
			State: importSource,
		},
		Schema: sch,
	}
}

func resourceHTTPSourceCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceSourceCreate(d, meta, func(source *api.SourceCreate) error {
		source.SourceType = "HTTP"
		raw, ok := d.GetOkExists("message_per_request")
		mpr := ok && raw.(bool)
		source.MessagePerRequest = &mpr
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

func resourceHTTPSourceRead(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceRead(d, meta, func(source *api.Source) error {
		d.Set("message_per_request", source.MessagePerRequest)
		return nil
	})
}

func resourceHTTPSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceUpdate(d, meta, func(source *api.Source) error {
		source.SourceType = "HTTP"
		raw, ok := d.GetOkExists("message_per_request")
		source.MessagePerRequest = ok && raw.(bool)
		return nil
	})
}
