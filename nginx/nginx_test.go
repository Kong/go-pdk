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
		{Method: "kong.nginx.get_var", Args: bridge.WrapString("foo"), Ret: bridge.WrapString("bar")},
	})

	ret, err := nginx.GetVar("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", ret)
}

func TestGetTLS1VersionStr(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{Method: "kong.nginx.get_tls1_version_str", Args: nil, Ret: bridge.WrapString("1.19")},
	})

	ret, err := nginx.GetTLS1VersionStr()
	assert.NoError(t, err)
	assert.Equal(t, "1.19", ret)
}

func TestCtx(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{Method: "kong.nginx.set_ctx", Args: &kong_plugin_protocol.KV{K: "key", V: structpb.NewStringValue("value")}, Ret: nil},
		{Method: "kong.nginx.get_ctx", Args: bridge.WrapString("key"), Ret: structpb.NewStringValue("value")},
		{Method: "kong.nginx.get_ctx", Args: bridge.WrapString("key_s"), Ret: structpb.NewStringValue("value")},
		{Method: "kong.nginx.get_ctx", Args: bridge.WrapString("key_n"), Ret: structpb.NewNumberValue(15.75)},
		{Method: "kong.nginx.get_ctx", Args: bridge.WrapString("key_i"), Ret: structpb.NewNumberValue(15)},
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
		{Method: "kong.nginx.req_start_time", Args: nil, Ret: &kong_plugin_protocol.Number{V: 1617060050.0}},
	})

	ret, err := nginx.ReqStartTime()
	assert.NoError(t, err)
	assert.Equal(t, 1617060050.0, ret)
}

func TestGetSubsystem(t *testing.T) {
	nginx := mockNginx(t, []bridgetest.MockStep{
		{Method: "kong.nginx.get_subsystem", Args: nil, Ret: bridge.WrapString("http")},
	})

	ret, err := nginx.GetSubsystem()
	assert.NoError(t, err)
	assert.Equal(t, "http", ret)
}
