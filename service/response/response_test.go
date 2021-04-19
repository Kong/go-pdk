package response

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	h, err := bridge.WrapHeaders(map[string][]string{
		"Host":         {"example.com"},
		"X-Two-Things": {"first", "second"},
	})
	assert.NoError(t, err)

	body := `GET / HTTP/1.1
Host: example.com
Accept: *

this is the content`

	response := Response{bridge.New(bridgetest.Mock(t, []bridgetest.MockStep{
		{Method: "kong.service.response.get_status", Args: nil, Ret: &kong_plugin_protocol.Int{V: 404}},
		{Method: "kong.service.response.get_header", Args: bridge.WrapString("Host"), Ret: bridge.WrapString("example.com")},
		{Method: "kong.service.response.get_headers", Args: &kong_plugin_protocol.Int{V: 30}, Ret: h},
		{Method: "kong.service.response.get_raw_body", Args: nil, Ret: bridge.WrapString(body)},
	}))}

	res_n, err := response.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, 404, res_n)

	res_s, err := response.GetHeader("Host")
	assert.NoError(t, err)
	assert.Equal(t, "example.com", res_s)

	res_h, err := response.GetHeaders(30)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"Host":         {"example.com"},
		"X-Two-Things": {"first", "second"},
	}, res_h)

	res_s, err = response.GetRawBody()
	assert.NoError(t, err)
	assert.Equal(t, body, res_s)
}
