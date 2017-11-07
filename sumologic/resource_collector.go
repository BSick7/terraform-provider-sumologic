package sumologic

import (
	"fmt"
	"strconv"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceCollectorCreate,
		Read:   resourceCollectorRead,
		Update: resourceCollectorUpdate,
		Delete: resourceCollectorDelete,
		Exists: resourceCollectorExists,

		Schema: map[string]*schema.Schema{
			"collector_type": {
				Type:     schema.TypeString,
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
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"collector_version": {
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
			"target_cpu": {
				Type:     schema.TypeFloat,
				Computed: true,
				Optional: true,
			},
			"alive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_seen_alive": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"cutoff_timestamp": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"cutoff_relative_time": {
				Type:     schema.TypeString,
				Computed: true,
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
			"source_sync_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCollectorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	newCollector := &api.CollectorCreate{
		CollectorType: d.Get("collector_type").(string),
		Name:          d.Get("name").(string),
	}

	collector, err := client.Collectors().Create(newCollector)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", collector.ID))

	if err := resourceCollectorUpdate(d, meta); err != nil {
		return err
	}

	return resourceCollectorRead(d, meta)
}

func resourceCollectorRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("collector_type", collector.CollectorType)
	d.Set("collector_version", collector.CollectorVersion)

	d.Set("host_name", collector.HostName)
	d.Set("time_zone", collector.TimeZone)
	d.Set("ephemeral", collector.Ephemeral)

	d.Set("alive", collector.Alive)
	d.Set("last_seen_alive", collector.LastSeenAlive)
	d.Set("cutoff_timestamp", collector.CutoffTimestamp)
	d.Set("cutoff_relative_time", collector.CutoffRelativeTime)

	d.Set("os_arch", collector.OsArch)
	d.Set("os_version", collector.OsVersion)
	d.Set("os_name", collector.OsName)
	d.Set("os_time", collector.OsTime)
	d.Set("target_cpu", collector.TargetCPU)
	d.Set("source_sync_mode", collector.SourceSyncMode)

	return nil
}

func resourceCollectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid id: %s", err)
	}

	collector := &api.Collector{ID: id}

	collector.Name = d.Get("name").(string)
	collector.Description = d.Get("description").(string)

	collector.HostName = d.Get("host_name").(string)
	collector.TimeZone = d.Get("time_zone").(string)
	collector.Ephemeral = d.Get("ephemeral").(bool)

	if raw, ok := d.GetOk("target_cpu"); ok {
		collector.TargetCPU = raw.(int64)
	}

	if raw, ok := d.GetOk("cutoff_timestamp"); ok {
		collector.CutoffTimestamp = raw.(int64)
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

	_, err = client.Collectors().Get(id)
	if aerr, ok := err.(*api.APIError); ok && aerr.Code == "InvalidCollector" {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
