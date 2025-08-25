package ranger

import (
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "url": {
        Type:        schema.TypeString,
        Required:    true,
        Description: "Base URL of the Ranger server (e.g. https://ranger-host:6080)",
      },
      "username": {
        Type:        schema.TypeString,
        Required:    true,
        Sensitive:   true,
      },
      "password": {
        Type:        schema.TypeString,
        Required:    true,
        Sensitive:   true,
      },
    },
    ResourcesMap: map[string]*schema.Resource{
      "ranger_policy": resourcePolicy(),
    },
    ConfigureContextFunc: providerConfigure,
  }
}
