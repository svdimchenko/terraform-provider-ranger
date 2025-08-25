package main

import (
  "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
  "github.com/svdimchenko/terraform-provider-ranger/ranger"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: ranger.Provider,
  })
}
