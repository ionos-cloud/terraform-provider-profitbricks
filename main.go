package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/ionos-cloud/terraform-provider-profitbricks/profitbricks"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: profitbricks.Provider})
}
