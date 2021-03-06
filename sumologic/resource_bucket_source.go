package sumologic

import (
	"time"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBucketSource(contentType string) *schema.Resource {
	return &schema.Resource{
		Create: resourceBucketSourceCreate(contentType),
		Read:   resourceBucketSourceRead,
		Update: resourceBucketSourceUpdate(contentType),
		Delete: resourceSourceDelete,
		Exists: resourceSourceExists,
		Importer: &schema.ResourceImporter{
			State: importSource,
		},
		Schema: resourceBucketSourceSchema(),
	}
}

func resourceBucketSourceSchema() map[string]*schema.Schema {
	sch := resourceSourceSchema()
	sch["paused"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	sch["scan_interval"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	sch["aws_bucket"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	sch["path_expression"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "*",
	}
	sch["aws_access_key"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	sch["aws_secret_key"] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	return sch
}

func resourceBucketSourceCreate(contentType string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		err := resourceSourceCreate(d, meta, func(source *api.SourceCreate) error {
			source.SourceType = "Polling"
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
		if err := resourceBucketSourceUpdate(contentType)(d, meta); err != nil {
			return err
		}
		return resourceBucketSourceRead(d, meta)
	}
}

func resourceBucketSourceRead(d *schema.ResourceData, meta interface{}) error {
	return resourceSourceRead(d, meta, func(source *api.Source) error {
		d.Set("paused", source.Paused)
		d.Set("scan_interval", source.ScanInterval.String())

		if tpr := source.ThirdPartyRef; tpr != nil && len(tpr.Resources) > 0 {
			if res := tpr.Resources[0]; res != nil {
				d.Set("aws_bucket", res.Path.BucketName)
				d.Set("path_expression", res.Path.PathExpression)
				d.Set("aws_access_key", res.Authentication.AccessKey)
			}
		}
		return nil
	})
}

func resourceBucketSourceUpdate(contentType string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		return resourceSourceUpdate(d, meta, func(source *api.Source) error {
			source.SourceType = "Polling"
			source.ContentType = contentType
			source.Paused = d.Get("paused").(bool)
			if raw, ok := d.GetOk("scan_interval"); ok {
				source.ScanInterval, _ = time.ParseDuration(raw.(string))
			}

			source.ThirdPartyRef = &api.ThirdPartyRef{
				Resources: []*api.ThirdPartyRefResource{
					{
						ServiceType: "AwsCloudFrontBucket",
						Path: &api.ThirdPartyRefResourcePath{
							Type:           "S3BucketPathExpression",
							BucketName:     d.Get("aws_bucket").(string),
							PathExpression: d.Get("path_expression").(string),
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
}
