package router

import (
	"github.com/kong/go-pdk/entities"
	"github.com/kong/go-pdk/bridge"
)

type Router struct {
	bridge.PdkBridge
}

func New(ch chan string) *Router {
	return &Router{*bridge.New(ch)}
}

func (c *Router) GetRoute() *entities.Route {
	reply := c.Ask(`kong.router.get_route`)
	if reply == "null" {
		return nil
	}
	route := entities.Route{}
	bridge.Unmarshal(reply, &route)
	return &route
}

func (c *Router) GetService() *entities.Service {
	reply := c.Ask(`kong.router.get_service`)
	if reply == "null" {
		return nil
	}
	service := entities.Service{}
	bridge.Unmarshal(reply, &service)
	return &service
}
