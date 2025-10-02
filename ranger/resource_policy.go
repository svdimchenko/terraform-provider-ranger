package ranger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parsePolicyResponse(respBody []byte, d *schema.ResourceData) (string, diag.Diagnostics) {
	var data map[string]interface{}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return "", diag.FromErr(err)
	}

	// Read the ID field as int
	if idValue, ok := data["id"]; ok {
		if id, ok := idValue.(float64); ok {
			d.SetId(fmt.Sprintf("%d", int(id)))
		}
	}

	// Remove redundant policy keys from terraform state
	delete(data, "id")
	delete(data, "guid")
	delete(data, "createdBy")
	delete(data, "updatedBy")
	delete(data, "createTime")
	delete(data, "updateTime")
	delete(data, "version")
	delete(data, "resourceSignature")

	// Convert back to JSON
	cleanedJSON, err := json.Marshal(data)
	if err != nil {
		return "", diag.FromErr(err)
	}
	return string(cleanedJSON), nil
}


func resourcePolicy() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourcePolicyCreate,
    ReadContext:   resourcePolicyRead,
    UpdateContext: resourcePolicyUpdate,
    DeleteContext: resourcePolicyDelete,

    Importer: &schema.ResourceImporter{
        StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
            client := m.(*Client)
            policyID := d.Id()

            // Call Ranger API to get policy by ID
            url := fmt.Sprintf("/service/public/v2/api/policy/%s", policyID)
            resp, err := client.rest.R().Get(url)
            if err != nil {
                return nil, fmt.Errorf("failed to import Ranger policy (ID: %s): %v", policyID, err)
            }
            if resp.StatusCode() == 404 {
                return nil, fmt.Errorf("policy with ID %s not found", policyID)
            }
            if resp.IsError() {
                return nil, fmt.Errorf("failed to import Ranger policy (ID: %s): %s", policyID, resp.String())
            }

            cleanedDef, diags := parsePolicyResponse(resp.Body(), d)
            if diags != nil {
                return nil, fmt.Errorf("failed to parse policy response: %v", diags)
            }
            d.Set("definition", cleanedDef)
            return []*schema.ResourceData{d}, nil
        },
    },

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
  // Parse response to get the policy ID
	var policy struct {
		ID int `json:"id"`
	}
  if err := json.Unmarshal(resp.Body(), &policy); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", policy.ID))
  time.Sleep(2 * time.Second)
  return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*Client)
  service := d.Get("service").(string)
  name := d.Get("name").(string)

  url := fmt.Sprintf("/service/public/v2/api/service/%s/policy/%s", service, name)
  resp, err := client.rest.R().Get(url)
  if resp.StatusCode() == 404 {
    return nil
  }
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.IsError() {
    return diag.Errorf("Failed to read policy: %s", resp.String())
  }

  cleanedDef, diags := parsePolicyResponse(resp.Body(), d)
  if diags != nil {
      return diags
  }
  d.Set("definition", cleanedDef)

  return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*Client)
  def := d.Get("definition").(string)
  policyID := d.Id()

  url := fmt.Sprintf("/service/public/v2/api/policy/%s", policyID)
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
  policyID := d.Id()

  url := fmt.Sprintf("/service/public/v2/api/policy/%s", policyID)
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
