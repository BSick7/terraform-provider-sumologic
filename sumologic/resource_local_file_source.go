package sumologic

import (
	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLocalFileSource() *schema.Resource {
	sch := resourceSourceSchema()
	sch["path_expression"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	sch["blacklist"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	sch["encoding"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "UTF-8",
	}

	return &schema.Resource{
		Create: resourceLocalFileSourceCreate,
		Read:   resourceLocalFileSourceRead,
		Update: resourceLocalFileSourceUpdate,
		Delete: resourceSourceDelete,
		Exists: resourceSourceExists,
		Importer: &schema.ResourceImporter{
			State: importSource,
		},
		Schema: sch,
	}
}

func resourceLocalFileSourceCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceSourceCreate(d, meta, func(source *api.SourceCreate) error {
		source.SourceType = "LocalFile"
		pe := d.Get("path_expression").(string)
		source.PathExpression = &pe
		return nil
	})
	if err != nil {
		return err
	}
	if err := resourceLocalFileSourceUpdate(d, meta); err != nil {
		return err
	}
	return resourceLocalFileSourceRead(d, meta)
}

func resourceLocalFileSourceRead(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceRead(d, meta, func(source *api.Source) error {
		d.Set("path_expression", source.PathExpression)
		d.Set("blacklist", source.Blacklist)
		d.Set("encoding", source.Encoding)
		return nil
	})
}

func resourceLocalFileSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceUpdate(d, meta, func(source *api.Source) error {
		source.SourceType = "LocalFile"
		if v := d.Get("blacklist").(*schema.Set); v.Len() > 0 {
			source.Blacklist = []string{}
			for _, bl := range v.List() {
				source.Blacklist = append(source.Blacklist, bl.(string))
			}
		}
		source.Encoding = d.Get("encoding").(string)
		return nil
	})
}
