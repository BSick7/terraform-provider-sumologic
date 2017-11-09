package sumologic

import (
	"fmt"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHTTPSource() *schema.Resource {
	sch := defaultSchema()
	sch["message_per_request"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Enable One Message Per Request",
	}

	return &schema.Resource{
		Create: resourceHTTPSourceCreate,
		Read:   readSourceFromSumologic,
		Update: resourceHTTPSourceUpdate,
		Delete: deleteSource,
		Exists: doesSourceExist,
		Importer: &schema.ResourceImporter{
			State: importSource,
		},
		Schema: sch,
	}
}

func resourceHTTPSourceCreate(d *schema.ResourceData, meta interface{}) error {
	source, err := createSource(d, meta, "HTTP")
	if err != nil {
		return err
	} else if source == nil {
		return fmt.Errorf("source was not created")
	}
	d.SetId(fmt.Sprintf("%d", source.ID))

	if err := resourceHTTPSourceUpdate(d, meta); err != nil {
		return err
	}

	return readSourceFromSumologic(d, meta)
}

func resourceHTTPSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	source, err := readSourceFromTerraform(d)
	if err != nil {
		return fmt.Errorf("error reading terraform values for http source: %s", err)
	}
	source.SourceType = "HTTP"

	if _, err := client.Collectors().Sources(collectorID).Update(source); err != nil {
		return err
	}
	return nil
}
