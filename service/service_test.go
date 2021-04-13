package service

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	service := Service{bridge.New(bridgetest.Mock(t, []bridgetest.MockStep{
		{"kong.service.set_upstream", bridge.WrapString("farm_4"), nil},
		{"kong.service.set_target", &kong_plugin_protocol.Target{Host: "internal.server.lan", Port: 8443}, nil},
	}))}

	assert.NoError(t, service.SetUpstream("farm_4"))
	assert.NoError(t, service.SetTarget("internal.server.lan", 8443))
}
