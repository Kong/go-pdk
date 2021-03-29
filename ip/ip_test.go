package ip

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)

func TestIsTrusted(t *testing.T) {
	ip := Ip{bridge.New(bridgetest.Mock(t, []bridgetest.MockStep{
		{"kong.ip.is_trusted", bridge.WrapString("1.1.1.1"), &kong_plugin_protocol.Bool{V: true}},
		{"kong.ip.is_trusted", bridge.WrapString("1.0.0.1"), &kong_plugin_protocol.Bool{V: false}},
	}))}

	ret, err := ip.IsTrusted("1.1.1.1")
	assert.NoError(t, err)
	assert.True(t, ret)

	ret, err = ip.IsTrusted("1.0.0.1")
	assert.NoError(t, err)
	assert.False(t, ret)
}

