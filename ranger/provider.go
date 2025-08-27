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
      "skip_tls_verify": {
        Type:        schema.TypeBool,
        Optional:    true,
        Sensitive:   false,
        Default: false,
        Description: "Skip TLS certificate verification (insecure, use only for testing)",
      },
    },
    ResourcesMap: map[string]*schema.Resource{
      "ranger_policy": resourcePolicy(),
    },
    ConfigureContextFunc: providerConfigure,
  }
}
