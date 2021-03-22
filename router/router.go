/*
Router module.

A set of functions to access the routing properties of the request.
*/
package router

import (
	"log"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
)

// Holds this module's functions.  Accessible as `kong.Router`
type Router struct {
	bridge.PdkBridge
}

// kong.Router.GetRoute() returns the current route entity.
// The request was matched against this route.
func (c Router) GetRoute() (route entities.Route, err error) {
	out := new(kong_plugin_protocol.Route)
	err = c.Ask(`kong.router.get_route`, nil, out)
	if err != nil {
		return
	}

	log.Printf("route: %v", out)

	route = entities.Route{
		Id:                      out.Id,
		CreatedAt:               int(out.CreatedAt),
		UpdatedAt:               int(out.UpdatedAt),
		Name:                    out.Name,
		Protocols:               out.Protocols,
		Methods:                 out.Methods,
		Hosts:                   out.Hosts,
		Paths:                   out.Paths,
		Headers:                 out.Headers,
		HTTPSRedirectStatusCode: int(out.HttpsRedirectStatusCode),
		RegexPriority:           int(out.RegexPriority),
		StripPath:               out.StripPath,
		PreserveHost:            out.PreserveHost,
		SNIs:                    out.Snis,
		Sources:                 out.Sources,
		Destinations:            out.Destinations,
		Tags:                    out.Tags,
		Service:                 entities.ServiceKey{Id: out.Service.Id},
	}
	return
}

// // kong.Router.GetService() returns the current service entity.
// // The request will be targetted to this upstream service.
func (c Router) GetService() (service entities.Service, err error) {
	out := new(kong_plugin_protocol.Service)
	err = c.Ask(`kong.router.get_service`, nil, out)
	if err != nil {
		return
	}

	service = entities.Service{
		Id:                out.Id,
		CreatedAt:         int(out.CreatedAt),
		UpdatedAt:         int(out.UpdatedAt),
		Name:              out.Name,
		Retries:           int(out.Retries),
		Protocol:          out.Protocol,
		Host:              out.Host,
		Port:              int(out.Port),
		Path:              out.Path,
		ConnectTimeout:    int(out.ConnectTimeout),
		WriteTimeout:      int(out.WriteTimeout),
		ReadTimeout:       int(out.ReadTimeout),
		Tags:              out.Tags,
		ClientCertificate: entities.CertificateKey{Id: out.Id},
	}
	return
}
