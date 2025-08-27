package ranger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// testSchema creates a schema.ResourceData with the provider fields.
func testSchema(t *testing.T) *schema.ResourceData {
	p := map[string]*schema.Schema{
		"url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"username": {
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:     schema.TypeString,
			Required: true,
		},
		"skip_tls_verify": {
			Type:     schema.TypeBool,
			Optional: true,
			Default: false,
		},
	}
	return schema.TestResourceDataRaw(t, p, nil)
}

func TestProviderConfigure_Success(t *testing.T) {
	// Mock Ranger API that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/service/public/v2/api/policy" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Build schema with test values
	d := testSchema(t)
	d.Set("url", server.URL)
	d.Set("username", "testuser")
	d.Set("password", "testpass")
	d.Set("skip_tls_verify", false)

	// Call configure
	clientIface, diags := providerConfigure(context.Background(), d)

	// Assert no diagnostics
	assert.Empty(t, diags)

	// Assert we got a valid client
	client, ok := clientIface.(*Client)
	assert.True(t, ok)
	assert.NotNil(t, client.rest)
	assert.Equal(t, server.URL, client.rest.BaseURL)
}

func TestProviderConfigure_Failure(t *testing.T) {
	// Server that always errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	d := testSchema(t)
	d.Set("url", server.URL)
	d.Set("username", "baduser")
	d.Set("password", "badpass")

	clientIface, diags := providerConfigure(context.Background(), d)

	// Assert diagnostics returned
	assert.Nil(t, clientIface)
	assert.NotEmpty(t, diags)
}
