/*
Service module.

The service module contains a set of functions to manipulate
the connection aspect of the request to the Service,
such as connecting to a given host, IP address/port,
or choosing a given Upstream entity for load-balancing and healthchecking.
*/
package service

import (
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.Service`
type Service struct {
	bridge.PdkBridge
}

// Called by the plugin server at initialization.
// func New(ch chan interface{}) Service {
// 	return Service{bridge.New(ch)}
// }

// kong.Service.SetUpstream() sets the desired Upstream entity to handle
// the load-balancing step for this request. Using this method is equivalent
// to creating a Service with a host property equal to that of an Upstream
// entity (in which case, the request would be proxied to one of the Targets
// associated with that Upstream).
//
// The host argument should receive a string equal to that of one of
// the Upstream entities currently configured.
func (s Service) SetUpstream(host string) error {
	return s.Ask(`kong.service.set_upstream`, bridge.WrapString(host), nil)
}

// kong.Service.SetTarget() sets the host and port on which to connect to
// for proxying the request.
//
// Using this method is equivalent to ask Kong to not run the load-balancing
// phase for this request, and consider it manually overridden.
// Load-balancing components such as retries and health-checks
// will also be ignored for this request.
//
// The host argument expects a string containing the IP address
// of the upstream server (IPv4/IPv6), and the port argument must
// contain a number representing the port on which to connect to.
func (s Service) SetTarget(host string, port int) error {
	arg := kong_plugin_protocol.Target{
		Host: host,
		Port: int32(port),
	}
	return s.Ask(`kong.service.set_target`, &arg, nil)
}

// TODO set_tls_cert_key
