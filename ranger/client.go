package ranger

import (
	"context"
	"crypto/tls"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Client struct {
  rest *resty.Client
}

func newClient(url, username, password string, skipTlsVerify bool) *Client {
  c := resty.New()
  c.SetBaseURL(url).
    SetBasicAuth(username, password).
    SetHeader("Content-Type", "application/json")
  if skipTlsVerify {
    c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
  }
  return &Client{rest: c}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
    url := d.Get("url").(string)
    username := d.Get("username").(string)
    password := d.Get("password").(string)
    skipTlsVerify := d.Get("skip_tls_verify").(bool)

    client := newClient(url, username, password, skipTlsVerify)

    resp, err := client.rest.R().Get("/service/public/v2/api/policy")
    if err != nil { return nil, diag.FromErr(err) }
    if resp.IsError() {
      return nil, diag.Errorf("failed to connect to Ranger API: %s", resp.Status())
    }

    return client, nil
}
