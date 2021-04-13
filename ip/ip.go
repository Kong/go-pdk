/*
Trusted IPs module.

This module can be used to determine whether or not a given IP address
is in the range of trusted IP addresses defined by the trusted_ips
configuration property.

Trusted IP addresses are those that are known to send correct replacement
addresses for clients (as per the chosen header field, e.g. X-Forwarded-*).

See https://docs.konghq.com/latest/configuration/#trusted_ips
*/
package ip

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
)

// Holds this module's functions.  Accessible as `kong.Ip`
type Ip struct {
	bridge.PdkBridge
}

// Depending on the trusted_ips configuration property, this function
// will return whether a given ip is trusted or not.
// Both ipv4 and ipv6 are supported.
func (ip Ip) IsTrusted(address string) (bool, error) {
	out := new(kong_plugin_protocol.Bool)
	err := ip.Ask(`kong.ip.is_trusted`, bridge.WrapString(address), out)
	if err != nil {
		return false, err
	}

	return out.V, nil
}
