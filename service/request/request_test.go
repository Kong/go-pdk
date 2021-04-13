package request

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestServRequest(t *testing.T) {
// 	q, err := bridge.WrapHeaders(map[string][]string{
// 		"ref":   []string{"wayback"},
// 		"trail": []string{"faint"},
// 	})
// 	assert.NoError(t, err)
//
// 	h, err := bridge.WrapHeaders(map[string][]string{
// 		"Host":         []string{"example.com"},
// 		"X-Two-Things": []string{"first", "second"},
// 	})
// 	assert.NoError(t, err)

	body := `GET / HTTP/1.1
Host: example.com
Accept: *

this is the content`

	request := Request{bridge.New(bridgetest.Mock(t, []bridgetest.MockStep{
		{"kong.service.request.set_scheme", bridge.WrapString("https"), nil},
		{"kong.service.request.set_path", bridge.WrapString("/login/orout"), nil},
		{"kong.service.request.set_raw_query", bridge.WrapString("ref=wayback&trail=faint"), nil},
		{"kong.service.request.set_method", bridge.WrapString("POST"), nil},
// 		{"kong.service.request.set_query", q, nil},
		{"kong.service.request.set_header", &kong_plugin_protocol.KV{K: "Host", V: structpb.NewStringValue("example.com")}, nil},
		{"kong.service.request.add_header", &kong_plugin_protocol.KV{K: "Host", V: structpb.NewStringValue("example.com")}, nil},
		{"kong.service.request.clear_header", bridge.WrapString("CORS"), nil},
// 		{"kong.service.request.set_headers", bridge.WrapString(""), nil},
		{"kong.service.request.set_raw_body", bridge.WrapString(body), nil},
	}))}

	assert.NoError(t, request.SetScheme("https"))
	assert.NoError(t, request.SetPath("/login/orout"))
	assert.NoError(t, request.SetRawQuery("ref=wayback&trail=faint"))
	assert.NoError(t, request.SetMethod("POST"))
// 	assert.NoError(t, request.SetQuery(map[string][]string{
// 		"ref":   []string{"wayback"},
// 		"trail": []string{"faint"},
// 	}))
	assert.NoError(t, request.SetHeader("Host", "example.com"))
	assert.NoError(t, request.AddHeader("Host", "example.com"))
	assert.NoError(t, request.ClearHeader("CORS"))
// 	assert.NoError(t, request.SetHeaders(map[string][]string{
// 		"Host":         []string{"example.com"},
// 		"X-Two-Things": []string{"first", "second"},
// 	}))
	assert.NoError(t, request.SetRawBody(body))

}
