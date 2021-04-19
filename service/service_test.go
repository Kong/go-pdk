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
		{Method: "kong.service.set_upstream", Args: bridge.WrapString("farm_4"), Ret: nil},
		{Method: "kong.service.set_target", Args: &kong_plugin_protocol.Target{Host: "internal.server.lan", Port: 8443}, Ret: nil},
	}))}

	assert.NoError(t, service.SetUpstream("farm_4"))
	assert.NoError(t, service.SetTarget("internal.server.lan", 8443))
}
