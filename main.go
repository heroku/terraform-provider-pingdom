package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/heroku/terraform-provider-pingdom/pingdom"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pingdom.Provider,
	})
}
