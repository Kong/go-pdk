package server

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchemaDict(t *testing.T) {
	type Config struct {
		JWKSURL             string `json:"jwks_url" kong:"required=true,default=https://example.com/.well-known/jwks.json"`
		CacheTTL            int    `json:"cache_ttl" kong:"default=3600"`
		AuthorizationHeader string `json:"authorization_header"`
		NoJsonTag           string `kong:"default=no_json_tag"`
	}

	schema := GetSchemaDict(reflect.TypeOf(Config{}))
	assert.Equal(t, schema, schemaDict{
		"type": "record",
		"fields": []schemaDict{
			{"jwks_url": schemaDict{"type": "string", "required": true, "default": "https://example.com/.well-known/jwks.json"}},
			{"cache_ttl": schemaDict{"type": "integer", "default": "3600"}},
			{"authorization_header": schemaDict{"type": "string"}},
			{"nojsontag": schemaDict{"type": "string", "default": "no_json_tag"}},
		},
	})
}
