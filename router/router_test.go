package router

import (
	"testing"

	"github.com/Kong/go-pdk/entities"
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)


func TestRouter(t *testing.T) {
	router := Router{bridge.New(bridgetest.Mock(t, []bridgetest.MockStep{
		{"kong.router.get_route", nil, &kong_plugin_protocol.Route{
			Id: "001:002",
			Name: "route_66",
			Protocols: []string{"http", "tcp"},
			Paths: []string{"/v0/left", "/v1/this"},
		}},
		{"kong.router.get_service", nil, &kong_plugin_protocol.Service{
			Id: "003:004",
			Name: "self_service",
			Protocol: "http",
			Path: "/v0/left",
		}},
	}))}

	ret_r, err := router.GetRoute()
	assert.NoError(t, err)
	assert.Equal(t, entities.Route{
		Id: "001:002",
		Name: "route_66",
		Protocols: []string{"http", "tcp"},
		Paths: []string{"/v0/left", "/v1/this"},
	}, ret_r)

	ret_s, err := router.GetService()
	assert.NoError(t, err)
	assert.Equal(t, entities.Service{
		Id: "003:004",
		Name: "self_service",
		Protocol: "http",
		Path: "/v0/left",
	}, ret_s)
}
