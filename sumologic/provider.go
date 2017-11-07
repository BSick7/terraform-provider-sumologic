package sumologic

import (
	"log"

	"github.com/BSick7/sumologic-sdk-go/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMO_ACCESS_ID", nil),
				Description: "The access ID for SumoLogic.",
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMO_ACCESS_KEY", nil),
				Description: "The access key for SumoLogic.",
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base address used to communicate with SumoLogic.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	session := api.DefaultSession()
	session.SetCredentials(d.Get("access_id").(string), d.Get("access_key").(string))

	if raw, ok := d.GetOk("address"); ok {
		log.Printf("[INFO] Setting SumoLogic address to %s\n", raw.(string))
		session.SetAddress(raw.(string))
	}

	log.Println("[INFO] Initializing SumoLogic client")
	client := api.NewClient(session)
	client.Discover()
	return client, nil
}
