package router

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
)

type Router struct {
	bridge.PdkBridge
}

func New(ch chan string) Router {
	return Router{bridge.New(ch)}
}

func (c Router) GetRoute() (*entities.Route, error) {
	reply, err := c.Ask(`kong.router.get_route`)
	if reply == "null" {
		return nil, err
	}
	route := entities.Route{}
	bridge.Unmarshal(reply, &route)
	return &route, nil
}

func (c Router) GetService() (*entities.Service, error) {
	reply, err := c.Ask(`kong.router.get_service`)
	if err != nil {
		return nil, err
	}
	service := entities.Service{}
	bridge.Unmarshal(reply, &service)
	return &service, nil
}
