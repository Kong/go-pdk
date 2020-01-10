/*
A set of functions to access the routing properties of the request.
*/
package router

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
)

type Router struct {
	bridge.PdkBridge
}

func New(ch chan interface{}) Router {
	return Router{bridge.New(ch)}
}

func (c Router) GetRoute() (route entities.Route, err error) {
	reply, err := c.Ask(`kong.router.get_route`)
	if err != nil {
		return
	}

	var ok bool
	if route, ok = reply.(entities.Route); !ok {
		err = bridge.ReturnTypeError("entities.Route")
	}
	return
}

func (c Router) GetService() (service entities.Service, err error) {
	val, err := c.Ask(`kong.router.get_service`)
	if err != nil {
		return
	}

	var ok bool
	if service, ok = val.(entities.Service); !ok {
		err = bridge.ReturnTypeError("entities.Service")
	}
	return
}
