package main

import (
	"github.com/BSick7/terraform-provider-sumologic/sumologic"
	"github.com/hashicorp/terraform/plugin"
)

var Version string

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
