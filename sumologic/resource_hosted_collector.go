package sumologic

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHostedCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostedCollectorCreate,
		Read:   resourceHostedCollectorRead,
		Update: resourceHostedCollectorUpdate,
		Delete: resourceHostedCollectorDelete,
		Exists: resourceHostedCollectorExists,
		Importer: &schema.ResourceImporter{
			State: importCollector,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
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
			"ephemeral": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"alive": {
				Type:     schema.TypeBool,
				Computed: true,
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

			"os_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceHostedCollectorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	newCollector := &api.CollectorCreate{
		CollectorType: "Hosted",
		Name:          d.Get("name").(string),
	}

	collector, err := client.Collectors().Create(newCollector)
	if err != nil {
		return err
	}
	if collector == nil {
		return fmt.Errorf("collector was not created")
	}
	d.SetId(fmt.Sprintf("%d", collector.ID))

	if err := resourceHostedCollectorUpdate(d, meta); err != nil {
		return err
	}

	return resourceHostedCollectorRead(d, meta)
}

func resourceHostedCollectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	collector, err := client.Collectors().Get(id)
	if err != nil {
		return err
	}

	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)

	d.Set("version", collector.CollectorVersion)

	d.Set("host_name", collector.HostName)
	d.Set("time_zone", collector.TimeZone)
	d.Set("ephemeral", collector.Ephemeral)

	d.Set("alive", collector.Alive)

	d.Set("os_arch", collector.OsArch)
	d.Set("os_version", collector.OsVersion)
	d.Set("os_name", collector.OsName)
	d.Set("os_time", collector.OsTime)

	// We store cutoff timestamp as string since tf doesn't support int64/time.Time
	d.Set("cutoff_timestamp", collector.CutoffTimestamp.Format(time.RFC3339))

	return nil
}

func resourceHostedCollectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	collector := &api.Collector{
		ID:            id,
		CollectorType: "Hosted",
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		HostName:      d.Get("host_name").(string),
		TimeZone:      d.Get("time_zone").(string),
		Ephemeral:     d.Get("ephemeral").(bool),
		Alive:         d.Get("alive").(bool),
	}

	// We store cutoff timestamp as string since tf doesn't support int64/time.Time
	if raw, ok := d.GetOk("cutoff_timestamp"); ok {
		collector.CutoffTimestamp, _ = time.Parse(time.RFC3339, raw.(string))
	}

	if _, err := client.Collectors().Update(collector); err != nil {
		return err
	}
	return nil
}

func resourceHostedCollectorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	return client.Collectors().Delete(&api.Collector{ID: id})
}

func resourceHostedCollectorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid id: %s", err)
	}

	_, err = client.Collectors().Get(id)
	if serr, ok := err.(*api.SumologicError); ok && serr.Status == 404 {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
