package response

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func mockResponse(t *testing.T, s []bridgetest.MockStep) Response {
	return Response{bridge.New(bridgetest.Mock(t, s))}
}

func TestResponse(t *testing.T) {
	h, err := bridge.WrapHeaders(map[string][]string{
		"Host":         {"example.com"},
		"X-Two-Things": {"first", "second"},
	})
	assert.NoError(t, err)

	response := mockResponse(t, []bridgetest.MockStep{
		{Method: "kong.response.get_status", Args: nil, Ret: &kong_plugin_protocol.Int{V: 404}},
		{Method: "kong.response.get_headers", Args: &kong_plugin_protocol.Int{V: 30}, Ret: h},
		{Method: "kong.response.get_source", Args: nil, Ret: bridge.WrapString("service")},
		{Method: "kong.response.set_status", Args: &kong_plugin_protocol.Int{V: 201}, Ret: nil},
		{Method: "kong.response.set_header", Args: &kong_plugin_protocol.KV{K: "key", V: structpb.NewStringValue("value")}, Ret: nil},
		{Method: "kong.response.add_header", Args: &kong_plugin_protocol.KV{K: "key", V: structpb.NewStringValue("value")}, Ret: nil},
		{Method: "kong.response.clear_header", Args: bridge.WrapString("key"), Ret: nil},
		{Method: "kong.response.set_headers", Args: nil, Ret: nil},
	})

	res_n, err := response.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, 404, res_n)

	res_h, err := response.GetHeaders(30)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"Host":         {"example.com"},
		"X-Two-Things": {"first", "second"},
	}, res_h)

	res_s, err := response.GetSource()
	assert.NoError(t, err)
	assert.Equal(t, "service", res_s)

	assert.NoError(t, response.SetStatus(201))

	assert.NoError(t, response.SetHeader("key", "value"))
	assert.NoError(t, response.AddHeader("key", "value"))
	assert.NoError(t, response.ClearHeader("key"))
	assert.NoError(t, response.SetHeaders(map[string][]string{
		"Host":         {"example.com"},
		"X-Two-Things": {"first", "second"},
	}))
}
