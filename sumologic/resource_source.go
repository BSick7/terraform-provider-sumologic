package sumologic

import (
	"fmt"
	"strconv"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSourceCreate,
		Read:   resourceSourceRead,
		Update: resourceSourceUpdate,
		Delete: resourceSourceDelete,
		Exists: resourceSourceExists,

		Schema: map[string]*schema.Schema{
			"collector_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"force_time_zone": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"automatic_date_parsing": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"multiline_processing_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"use_autoline_matching": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"manual_prefix_regexp": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"default_date_format": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"default_date_formats": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Optional: true,
			},
			"filters": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Optional: true,
			},
			"cutoff_timestamp": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"cutoff_relative_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceSourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	newSource := &api.Source{
		SourceType: d.Get("type").(string),
		Name:       d.Get("name").(string),
	}

	source, err := client.Collectors().Sources(collectorID).Create(newSource)
	if err != nil {
		return err
	}
	if source == nil {
		return fmt.Errorf("source was not created")
	}
	d.SetId(fmt.Sprintf("%d", source.ID))

	if err := resourceSourceUpdate(d, meta); err != nil {
		return err
	}

	return resourceSourceRead(d, meta)
}

func resourceSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	source, err := client.Collectors().Sources(collectorID).Get(id)
	if err != nil {
		return err
	}

	d.Set("name", source.Name)
	d.Set("type", source.SourceType)
	d.Set("description", source.Description)
	d.Set("category", source.Category)

	d.Set("host_name", source.HostName)
	d.Set("time_zone", source.TimeZone)
	d.Set("force_time_zone", source.ForceTimeZone)

	d.Set("automatic_date_parsing", source.AutomaticDateParsing)
	d.Set("multiline_processing_enabled", source.MultilineProcessingEnabled)
	d.Set("use_autoline_matching", source.UseAutolineMatching)
	d.Set("manual_prefix_regexp", source.ManualPrefixRegexp)

	d.Set("default_date_format", source.DefaultDateFormat)
	d.Set("default_date_formats", source.DefaultDateFormats)

	d.Set("filters", source.Filters)

	d.Set("cutoff_timestamp", source.CutoffTimestamp)
	d.Set("cutoff_relative_time", source.CutoffRelativeTime)

	return nil
}

func resourceSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	source := &api.Source{ID: id}

	source.Name = d.Get("name").(string)
	source.SourceType = d.Get("type").(string)
	source.Description = d.Get("description").(string)
	source.Category = d.Get("category").(string)

	source.HostName = d.Get("host_name").(string)
	source.TimeZone = d.Get("time_zone").(string)
	source.ForceTimeZone = d.Get("force_time_zone").(bool)

	source.AutomaticDateParsing = d.Get("automatic_date_parsing").(bool)
	source.MultilineProcessingEnabled = d.Get("multiline_processing_enabled").(bool)
	source.UseAutolineMatching = d.Get("use_autoline_matching").(bool)

	source.ManualPrefixRegexp = d.Get("manual_prefix_regexp").(string)
	source.DefaultDateFormat = d.Get("default_date_format").(string)

	// TODO: source.DefaultDateFormats "default_date_formats"
	// TODO: source.Filters "filters"

	if raw, ok := d.GetOk("cutoff_timestamp"); ok {
		source.CutoffTimestamp = raw.(int64)
	}

	if _, err := client.Collectors().Sources(collectorID).Update(source); err != nil {
		return err
	}
	return nil
}

func resourceSourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	return client.Collectors().Sources(collectorID).Delete(&api.Source{ID: id})
}

func resourceSourceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid id: %s", err)
	}

	_, err = client.Collectors().Sources(collectorID).Get(id)
	if serr, ok := err.(*api.SumologicError); ok && serr.Status == 404 {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
