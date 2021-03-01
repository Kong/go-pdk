/*
Router module.

A set of functions to access the routing properties of the request.
*/
package router

import (
	"github.com/Kong/go-pdk/bridge"
// 	"github.com/Kong/go-pdk/entities"
)

// Holds this module's functions.  Accessible as `kong.Router`
type Router struct {
	bridge.PdkBridge
}

// Called by the plugin server at initialization.
// func New(ch chan interface{}) Router {
// 	return Router{bridge.New(ch)}
// }

// kong.Router.GetRoute() returns the current route entity.
// The request was matched against this route.
// func (c Router) GetRoute() (route entities.Route, err error) {
// 	reply, err := c.Ask(`kong.router.get_route`)
// 	if err != nil {
// 		return
// 	}
//
// 	var ok bool
// 	if route, ok = reply.(entities.Route); !ok {
// 		err = bridge.ReturnTypeError("entities.Route")
// 	}
// 	return
// }
//
// // kong.Router.GetService() returns the current service entity.
// // The request will be targetted to this upstream service.
// func (c Router) GetService() (service entities.Service, err error) {
// 	val, err := c.Ask(`kong.router.get_service`)
// 	if err != nil {
// 		return
// 	}
//
// 	var ok bool
// 	if service, ok = val.(entities.Service); !ok {
// 		err = bridge.ReturnTypeError("entities.Service")
// 	}
// 	return
// }
