package request

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)

func mockRequest(t *testing.T, s []bridgetest.MockStep) Request {
	return Request{bridge.New(bridgetest.Mock(t, s))}
}

func TestGetInfos(t *testing.T) {

	q, err := bridge.WrapHeaders(map[string][]string{
		"ref":   []string{"wayback"},
		"trail": []string{"faint"},
	})
	assert.NoError(t, err)

	h, err := bridge.WrapHeaders(map[string][]string{
		"Host":         []string{"example.com"},
		"X-Two-Things": []string{"first", "second"},
	})
	assert.NoError(t, err)

	body := `GET / HTTP/1.1
Host: example.com
Accept: *

this is the content`

	request := mockRequest(t, []bridgetest.MockStep{
		{"kong.request.get_scheme", nil, bridge.WrapString("https")},
		{"kong.request.get_host", nil, bridge.WrapString("example.com")},
		{"kong.request.get_port", nil, &kong_plugin_protocol.Int{V: 443}},
		{"kong.request.get_forwarded_scheme", nil, bridge.WrapString("https")},
		{"kong.request.get_forwarded_host", nil, bridge.WrapString("example.com")},
		{"kong.request.get_forwarded_port", nil, &kong_plugin_protocol.Int{V: 443}},
		{"kong.request.get_http_version", nil, &kong_plugin_protocol.Number{V: 2.1}},
		{"kong.request.get_method", nil, bridge.WrapString("HEADER")},
		{"kong.request.get_path", nil, bridge.WrapString("/login/orout")},
		{"kong.request.get_path_with_query", nil, bridge.WrapString("/login/orout?ref=wayback")},
		{"kong.request.get_raw_query", nil, bridge.WrapString("ref=wayback&trail=faint")},
		{"kong.request.get_query_arg", bridge.WrapString("ref"), bridge.WrapString("wayback")},
		{"kong.request.get_query", &kong_plugin_protocol.Int{V: 30}, q},
		{"kong.request.get_header", bridge.WrapString("Host"), bridge.WrapString("example.com")},
		{"kong.request.get_headers", &kong_plugin_protocol.Int{V: 30}, h},
		{"kong.request.get_raw_body", nil, bridge.WrapString(body)},
	})

	ret, err := request.GetScheme()
	assert.NoError(t, err)
	assert.Equal(t, "https", ret)

	ret, err = request.GetHost()
	assert.NoError(t, err)
	assert.Equal(t, "example.com", ret)

	ret_i, err := request.GetPort()
	assert.NoError(t, err)
	assert.Equal(t, 443, ret_i)

	ret, err = request.GetForwardedScheme()
	assert.NoError(t, err)
	assert.Equal(t, "https", ret)

	ret, err = request.GetForwardedHost()
	assert.NoError(t, err)
	assert.Equal(t, "example.com", ret)

	ret_i, err = request.GetForwardedPort()
	assert.NoError(t, err)
	assert.Equal(t, 443, ret_i)

	ret_f, err := request.GetHttpVersion()
	assert.NoError(t, err)
	assert.Equal(t, 2.1, ret_f)

	ret, err = request.GetMethod()
	assert.NoError(t, err)
	assert.Equal(t, "HEADER", ret)

	ret, err = request.GetPath()
	assert.NoError(t, err)
	assert.Equal(t, "/login/orout", ret)

	ret, err = request.GetPathWithQuery()
	assert.NoError(t, err)
	assert.Equal(t, "/login/orout?ref=wayback", ret)

	ret, err = request.GetRawQuery()
	assert.NoError(t, err)
	assert.Equal(t, "ref=wayback&trail=faint", ret)

	ret, err = request.GetQueryArg("ref")
	assert.NoError(t, err)
	assert.Equal(t, "wayback", ret)

	ret_q, err := request.GetQuery(30)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"ref":   []string{"wayback"},
		"trail": []string{"faint"},
	}, ret_q)

	ret, err = request.GetHeader("Host")
	assert.NoError(t, err)
	assert.Equal(t, "example.com", ret)

	ret_q, err = request.GetHeaders(30)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"Host":         []string{"example.com"},
		"X-Two-Things": []string{"first", "second"},
	}, ret_q)

	ret, err = request.GetRawBody()
	assert.NoError(t, err)
	assert.Equal(t, body, ret)
}
