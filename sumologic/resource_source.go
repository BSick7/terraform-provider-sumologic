package sumologic

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

// All Source Types and parameters:
// https://help.sumologic.com/Send-Data/Sources/03Use-JSON-to-Configure-Sources

func defaultSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"collector_id": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"host": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"category": {
			Type:     schema.TypeString,
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
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"format": {
						Type:     schema.TypeString,
						Required: true,
					},
					"locator": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"filters": {
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"regexp": {
						Type:     schema.TypeString,
						Required: true,
					},
					"mask": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"cutoff_timestamp": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Only collect data more recent than this timestamp (RFC3339 Formatted)",
		},
		"cutoff_relative_time": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
			Description: `Can be specified instead of cutoffTimestamp to provide a relative offset with respect to the current time.
Example: use -1h, -1d, or -1w to collect data that's less than one hour, one day, or one week old, respectively.
You can only use hours, days, and weeks to specify cutoffRelativeTime. No other time units are supported.`,
		},
	}
}

func readSourceFromTerraform(d *schema.ResourceData) (*api.Source, error) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid source id: %s", err)
	}

	source := &api.Source{
		ID:                         id,
		Name:                       d.Get("name").(string),
		Description:                d.Get("description").(string),
		Category:                   d.Get("category").(string),
		HostName:                   d.Get("host").(string),
		TimeZone:                   d.Get("time_zone").(string),
		ForceTimeZone:              d.Get("force_time_zone").(bool),
		AutomaticDateParsing:       d.Get("automatic_date_parsing").(bool),
		MultilineProcessingEnabled: d.Get("multiline_processing_enabled").(bool),
		UseAutolineMatching:        d.Get("use_autoline_matching").(bool),
		ManualPrefixRegexp:         d.Get("manual_prefix_regexp").(string),
		DefaultDateFormat:          d.Get("default_date_format").(string),
		CutoffRelativeTime:         d.Get("cutoff_relative_time").(string),
	}

	// We store cutoff timestamp as string since tf doesn't support int64/time.Time
	if raw, ok := d.GetOk("cutoff_timestamp"); ok {
		source.CutoffTimestamp, _ = time.Parse(time.RFC3339, raw.(string))
	}

	source.DefaultDateFormats = readSourceDefaultDateFormatsFromTerraform(d)
	source.Filters = readSourceFiltersFromTerraform(d)

	return source, nil
}

func readSourceDefaultDateFormatsFromTerraform(d *schema.ResourceData) []*api.DateFormat {
	formats := make([]*api.DateFormat, 0)

	if v, ok := d.GetOk("default_date_formats"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			sddf := v.(map[string]interface{})
			format := &api.DateFormat{
				Format:  sddf["format"].(string),
				Locator: sddf["locator"].(string),
			}
			formats = append(formats, format)
		}
	}

	return formats
}

func readSourceFiltersFromTerraform(d *schema.ResourceData) []*api.SourceFilter {
	filters := make([]*api.SourceFilter, 0)

	if v, ok := d.GetOk("filters"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			sf := v.(map[string]interface{})
			filter := &api.SourceFilter{
				FilterType: sf["type"].(string),
				Name:       sf["name"].(string),
				Regexp:     sf["regexp"].(string),
				Mask:       sf["mask"].(string),
			}
			filters = append(filters, filter)
		}
	}

	return filters
}

func readSourceFromSumologic(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	source, err := client.Collectors().Sources(collectorID).Get(id)
	if err != nil {
		return err
	} else if source == nil {
		d.SetId("")
		return nil
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

	d.Set("cutoff_relative_time", source.CutoffRelativeTime)

	// We store cutoff timestamp as string since tf doesn't support int64/time.Time
	d.Set("cutoff_timestamp", source.CutoffTimestamp.Format(time.RFC3339))

	return nil
}

func createSource(d *schema.ResourceData, meta interface{}, sourceType string) (*api.Source, error) {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	mpr := d.Get("message_per_request").(bool)

	newSource := &api.SourceCreate{
		SourceType:        sourceType,
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Category:          d.Get("category").(string),
		MessagePerRequest: &mpr,
	}

	return client.Collectors().Sources(collectorID).Create(newSource)
}

func deleteSource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	return client.Collectors().Sources(collectorID).Delete(&api.Source{ID: id})
}

func doesSourceExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*api.Client)
	collectorID := d.Get("collector_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid id: %s", err)
	}

	return client.Collectors().Sources(collectorID).Exists(id)
}
