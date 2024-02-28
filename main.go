package main

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: byteplus.Provider,
	})
}
