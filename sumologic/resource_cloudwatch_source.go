package sumologic

import (
	"time"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudWatchSource() *schema.Resource {
	sch := resourceBucketSourceSchema()
	delete(sch, "bucket_name")
	delete(sch, "path_expression")
	sch["limit_to_regions"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      schema.HashString,
	}
	sch["limit_to_namespaces"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      schema.HashString,
	}

	return &schema.Resource{
		Create: resourceCloudWatchSourceCreate,
		Read:   resourceCloudWatchSourceRead,
		Update: resourceCloudWatchSourceUpdate,
		Delete: resourceSourceDelete,
		Exists: resourceSourceExists,
		Importer: &schema.ResourceImporter{
			State: importSource,
		},
		Schema: sch,
	}
}

func resourceCloudWatchSourceCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceSourceCreate(d, meta, func(source *api.SourceCreate) error {
		source.SourceType = "Polling"
		contentType := "AwsCloudWatch"
		source.ContentType = &contentType
		paused := d.Get("paused").(bool)
		source.Paused = &paused
		if raw, ok := d.GetOk("scan_interval"); ok {
			si, _ := time.ParseDuration(raw.(string))
			source.ScanInterval = &si
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := resourceCloudWatchSourceUpdate(d, meta); err != nil {
		return err
	}
	return resourceCloudWatchSourceRead(d, meta)
}

func resourceCloudWatchSourceRead(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceRead(d, meta, func(source *api.Source) error {
		d.Set("paused", source.Paused)
		d.Set("scan_interval", source.ScanInterval.String())

		if tpr := source.ThirdPartyRef; tpr != nil && len(tpr.Resources) > 0 {
			if res := tpr.Resources[0]; res != nil {
				d.Set("limit_to_regions", res.Path.LimitToRegions)
				d.Set("limit_to_namespaces", res.Path.LimitToNamespaces)
				d.Set("aws_access_key", res.Authentication.AccessKey)
			}
		}
		return nil
	})
}

func resourceCloudWatchSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceUpdate(d, meta, func(source *api.Source) error {
		source.SourceType = "Polling"
		source.ContentType = "AwsCloudWatch"
		source.Paused = d.Get("paused").(bool)
		if raw, ok := d.GetOk("scan_interval"); ok {
			source.ScanInterval, _ = time.ParseDuration(raw.(string))
		}

		source.ThirdPartyRef = &api.ThirdPartyRef{
			Resources: []*api.ThirdPartyRefResource{
				{
					ServiceType: "AwsCloudWatch",
					Path: &api.ThirdPartyRefResourcePath{
						Type:              "S3BucketPathExpression",
						LimitToRegions:    readStringSliceFromTerraform(d, "limit_to_regions"),
						LimitToNamespaces: readStringSliceFromTerraform(d, "limit_to_namespaces"),
					},
					Authentication: &api.ThirdPartyRefResourceAuthentication{
						Type:      "S3BucketAuthentication",
						AccessKey: d.Get("aws_access_key").(string),
						SecretKey: d.Get("aws_secret_key").(string),
					},
				},
			},
		}

		return nil
	})
}

func readStringSliceFromTerraform(d *schema.ResourceData, key string) []string {
	v, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	var items []string
	for _, item := range v.(*schema.Set).List() {
		items = append(items, item.(string))
	}
	return items
}
