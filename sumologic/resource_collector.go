package sumologic

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCollectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func resourceCollectorCreate(d *schema.ResourceData, meta interface{}, custom func(*api.CollectorCreate) error) error {
	client := meta.(*api.Client)

	newCollector := &api.CollectorCreate{
		Name: d.Get("name").(string),
	}

	if custom != nil {
		if err := custom(newCollector); err != nil {
			return err
		}
	}

	collector, err := client.Collectors().Create(newCollector)
	if err != nil {
		return err
	} else if collector == nil {
		return fmt.Errorf("collector was not created")
	}
	d.SetId(fmt.Sprintf("%d", collector.ID))
	return nil
}

func resourceCollectorRead(d *schema.ResourceData, meta interface{}, custom func(collector *api.Collector) error) error {
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

	if custom != nil {
		if err := custom(collector); err != nil {
			return err
		}
	}

	return nil
}

func resourceCollectorUpdate(d *schema.ResourceData, meta interface{}, custom func(collector *api.Collector) error) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	collector := &api.Collector{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Category:    d.Get("category").(string),
		HostName:    d.Get("host_name").(string),
		TimeZone:    d.Get("time_zone").(string),
		Ephemeral:   d.Get("ephemeral").(bool),
		Alive:       d.Get("alive").(bool),
	}

	// We store cutoff timestamp as string since tf doesn't support int64/time.Time
	if raw, ok := d.GetOk("cutoff_timestamp"); ok {
		collector.CutoffTimestamp, _ = time.Parse(time.RFC3339, raw.(string))
	}

	if custom != nil {
		if err := custom(collector); err != nil {
			return err
		}
	}

	if _, err := client.Collectors().Update(collector); err != nil {
		return err
	}
	return nil
}

func resourceCollectorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	return client.Collectors().Delete(&api.Collector{ID: id})
}

func resourceCollectorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid id: %s", err)
	}

	return client.Collectors().Exists(id)
}
