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
		{Method: "kong.service.request.set_scheme", Args: bridge.WrapString("https"), Ret: nil},
		{Method: "kong.service.request.set_path", Args: bridge.WrapString("/login/orout"), Ret: nil},
		{Method: "kong.service.request.set_raw_query", Args: bridge.WrapString("ref=wayback&trail=faint"), Ret: nil},
		{Method: "kong.service.request.set_method", Args: bridge.WrapString("POST"), Ret: nil},
		// 		{Method:"kong.service.request.set_query",Args: q,Ret:  nil},
		{Method: "kong.service.request.set_header", Args: &kong_plugin_protocol.KV{K: "Host", V: structpb.NewStringValue("example.com")}, Ret: nil},
		{Method: "kong.service.request.add_header", Args: &kong_plugin_protocol.KV{K: "Host", V: structpb.NewStringValue("example.com")}, Ret: nil},
		{Method: "kong.service.request.clear_header", Args: bridge.WrapString("CORS"), Ret: nil},
		// 		{Method:"kong.service.request.set_headers",Args: bridge.WrapString(""),Ret:  nil},
		{Method: "kong.service.request.set_raw_body", Args: bridge.WrapString(body), Ret: nil},
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
