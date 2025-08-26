package ranger

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestProviderSchema(t *testing.T) {
	p := Provider()
	assert.NotNil(t, p)

	// check schema keys
	expectedKeys := []string{"url", "username", "password"}
	for _, key := range expectedKeys {
		_, ok := p.Schema[key]
		assert.Truef(t, ok, "expected key %s in schema", key)
	}

	// validate attributes
	urlSchema := p.Schema["url"]
	assert.Equal(t, schema.TypeString, urlSchema.Type)
	assert.True(t, urlSchema.Required)

	usernameSchema := p.Schema["username"]
	assert.Equal(t, schema.TypeString, usernameSchema.Type)
	assert.True(t, usernameSchema.Required)
	assert.True(t, usernameSchema.Sensitive)

	passwordSchema := p.Schema["password"]
	assert.Equal(t, schema.TypeString, passwordSchema.Type)
	assert.True(t, passwordSchema.Required)
	assert.True(t, passwordSchema.Sensitive)
}

func TestProviderResources(t *testing.T) {
	p := Provider()
	assert.NotNil(t, p)

	// check resources
	_, ok := p.ResourcesMap["ranger_policy"]
	assert.True(t, ok, "expected ranger_policy resource in ResourcesMap")
}

func TestProviderValidate(t *testing.T) {
	p := Provider()

	err := p.InternalValidate()
	assert.NoError(t, err, "provider validation should not return error")
}
