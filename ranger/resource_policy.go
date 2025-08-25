package ranger

import (
  "context"
  "fmt"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourcePolicyCreate,
    ReadContext:   resourcePolicyRead,
    UpdateContext: resourcePolicyUpdate,
    DeleteContext: resourcePolicyDelete,

    Schema: map[string]*schema.Schema{
      "service": {
        Type:     schema.TypeString,
        Required: true,
      },
      "name": {
        Type:     schema.TypeString,
        Required: true,
      },
      "definition": {
        Type:     schema.TypeString,
        Required: true,
        Description: "Full JSON body of the Ranger policy definition",
      },
    },
  }
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*Client)
  service := d.Get("service").(string)
  name := d.Get("name").(string)
  def := d.Get("definition").(string)

  resp, err := client.rest.R().
    SetBody(def).
    Post("/service/public/v2/api/policy")
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.IsError() {
    return diag.Errorf("Failed to create policy: %s", resp.String())
  }

  d.SetId(fmt.Sprintf("%s/%s", service, name))
  return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*Client)
  service := d.Get("service").(string)
  name := d.Get("name").(string)

  url := fmt.Sprintf("/service/public/v2/api/service/%s/policy/%s", service, name)
  resp, err := client.rest.R().Get(url)
  if resp.StatusCode() == 404 {
    d.SetId("")
    return nil
  }
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.IsError() {
    return diag.Errorf("Failed to read policy: %s", resp.String())
  }

  d.Set("definition", string(resp.Body()))
  return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*Client)
  def := d.Get("definition").(string)
  name := d.Get("name").(string)

  url := fmt.Sprintf("/service/public/v2/api/policy/%s", name)
  resp, err := client.rest.R().SetBody(def).Put(url)
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.IsError() {
    return diag.Errorf("Failed to update policy: %s", resp.String())
  }
  return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*Client)
  service := d.Get("service").(string)
  name := d.Get("name").(string)

  url := fmt.Sprintf("/service/public/v2/api/service/%s/policy/%s", service, name)
  resp, err := client.rest.R().Delete(url)
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.IsError() && resp.StatusCode() != 404 {
    return diag.Errorf("Failed to delete policy: %s", resp.String())
  }

  d.SetId("")
  return nil
}
