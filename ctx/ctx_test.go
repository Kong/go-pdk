package ctx

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func mockCtx(t *testing.T, s []bridgetest.MockStep) Ctx {
	return Ctx{bridge.New(bridgetest.Mock(t, s))}
}

func TestSetShared(t *testing.T) {
	ctx := mockCtx(t, []bridgetest.MockStep{
		{Method: "kong.ctx.shared.set", Args: &kong_plugin_protocol.KV{K: "key", V: structpb.NewStringValue("value")}, Ret: nil},
	})

	assert.NoError(t, ctx.SetShared("key", "value"))
}

func TestGetSharedAny(t *testing.T) {
	v, err := structpb.NewValue(67)
	assert.NoError(t, err)

	ctx := mockCtx(t, []bridgetest.MockStep{
		{Method: "kong.ctx.shared.get", Args: bridge.WrapString("key"), Ret: v},
	})

	val, err := ctx.GetSharedAny("key")
	assert.NoError(t, err)
	assert.Equal(t, 67.0, val)
}

func TestGetSharedString(t *testing.T) {
	ctx := mockCtx(t, []bridgetest.MockStep{
		{Method: "kong.ctx.shared.get", Args: bridge.WrapString("key"), Ret: structpb.NewStringValue("value")},
	})

	val, err := ctx.GetSharedString("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
}

func TestGetSharedFloat(t *testing.T) {
	ctx := mockCtx(t, []bridgetest.MockStep{
		{Method: "kong.ctx.shared.get", Args: bridge.WrapString("key"), Ret: structpb.NewNumberValue(2.74)},
	})

	val, err := ctx.GetSharedFloat("key")
	assert.NoError(t, err)
	assert.Equal(t, 2.74, val)
}
