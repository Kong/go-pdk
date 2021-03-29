package nginx

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func mockNginx(t *testing.T, s []bridgetest.MockStep) Nginx {
	return Nginx{bridge.New(bridgetest.Mock(t, s))}
}

func TestGetVar(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{"kong.nginx.get_var", bridge.WrapString("foo"), bridge.WrapString("bar")},
	})

	ret, err := nginx.GetVar("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", ret)
}

func TestGetTLS1VersionStr(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{"kong.nginx.get_tls1_version_str", nil, bridge.WrapString("1.19")},
	})

	ret, err := nginx.GetTLS1VersionStr()
	assert.NoError(t, err)
	assert.Equal(t, "1.19", ret)
}

func TestCtx(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{"kong.nginx.set_ctx", &kong_plugin_protocol.KV{K: "key", V: structpb.NewStringValue("value")}, nil},
		{"kong.nginx.get_ctx", bridge.WrapString("key"), structpb.NewStringValue("value")},
		{"kong.nginx.get_ctx", bridge.WrapString("key_s"), structpb.NewStringValue("value")},
		{"kong.nginx.get_ctx", bridge.WrapString("key_n"), structpb.NewNumberValue(15.75)},
		{"kong.nginx.get_ctx", bridge.WrapString("key_i"), structpb.NewNumberValue(15)},
	})

	assert.NoError(t, nginx.SetCtx("key", "value"))

	ret, err := nginx.GetCtxAny("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", ret)

	ret_s, err := nginx.GetCtxString("key_s")
	assert.NoError(t, err)
	assert.Equal(t, "value", ret_s)

	ret_n, err := nginx.GetCtxFloat("key_n")
	assert.NoError(t, err)
	assert.Equal(t, 15.75, ret_n)

	ret_i, err := nginx.GetCtxInt("key_i")
	assert.NoError(t, err)
	assert.Equal(t, 15, ret_i)
}

func TestReqStartTime(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{"kong.nginx.req_start_time", nil, &kong_plugin_protocol.Number{V: 1617060050.0}},
	})

	ret, err := nginx.ReqStartTime()
	assert.NoError(t, err)
	assert.Equal(t, 1617060050.0, ret)
}

func TestGetSubsystem(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{"kong.nginx.get_subsystem", nil, bridge.WrapString("http")},
	})

	ret, err := nginx.GetSubsystem()
	assert.NoError(t, err)
	assert.Equal(t, "http", ret)
}
