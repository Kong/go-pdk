/*
Package Kong/go-pdk implements Kong's Plugin Development Kit for Go.

It directly parallels the existing kong PDK for Lua plugins.

Kong plugins written in Go implement event handlers as methods on the Plugin's
structure, with the given signature:

	func (conf *MyConfig) Access (kong *pdk.PDK) {
		...
	}

The `kong` argument of type `*pdk.PDK` is the entrypoint for all PDK functions.
For example, to get the client's IP address, you'd use `kong.Client.GetIp()`.
*/
package pdk

import (
	"net"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/client"
	"github.com/Kong/go-pdk/ctx"
	"github.com/Kong/go-pdk/ip"
	"github.com/Kong/go-pdk/log"
	"github.com/Kong/go-pdk/nginx"
	"github.com/Kong/go-pdk/node"
	"github.com/Kong/go-pdk/request"
	"github.com/Kong/go-pdk/response"
	"github.com/Kong/go-pdk/router"
	"github.com/Kong/go-pdk/service"
	service_request "github.com/Kong/go-pdk/service/request"
	service_response "github.com/Kong/go-pdk/service/response"
)

// PDK go pdk module
type PDK struct {
	Client          client.Client
	Ctx             ctx.Ctx
	Log             log.Log
	Nginx           nginx.Nginx
	Request         request.Request
	Response        response.Response
	Router          router.Router
	IP              ip.Ip
	Node            node.Node
	Service         service.Service
	ServiceRequest  service_request.Request
	ServiceResponse service_response.Response
}

// Init initialize go pdk.  Called by the pluginserver at initialization.
func Init(conn net.Conn) *PDK {
	b := bridge.New(conn)
	return &PDK{
		Client:          client.Client{PdkBridge: b},
		Ctx:             ctx.Ctx{PdkBridge: b},
		Log:             log.Log{PdkBridge: b},
		Nginx:           nginx.Nginx{PdkBridge: b},
		Request:         request.Request{PdkBridge: b},
		Response:        response.Response{PdkBridge: b},
		Router:          router.Router{PdkBridge: b},
		IP:              ip.Ip{PdkBridge: b},
		Node:            node.Node{PdkBridge: b},
		Service:         service.Service{PdkBridge: b},
		ServiceRequest:  service_request.Request{PdkBridge: b},
		ServiceResponse: service_response.Response{PdkBridge: b},
	}
}
