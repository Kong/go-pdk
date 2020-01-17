/*
Client information module.

A set of functions to retrieve information about the client
connecting to Kong in the context of a given request.
*/
package client

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
	"fmt"
)

// Holds this module's functions.  Accessible as `kong.Client`
type Client struct {
	bridge.PdkBridge
}

func checkConsumer(v interface{}) (consumer entities.Consumer, err error) {
	consumer, ok := v.(entities.Consumer)
	if !ok {
		err = bridge.ReturnTypeError("Consumer Entity")
	}
	return
}

type AuthenticatedCredential struct {
	Id         string `json:"id"`
	ConsumerId string `json:"consumer_id"`
}

// Called by the plugin server at initialization.
func New(ch chan interface{}) Client {
	return Client{bridge.New(ch)}
}

// kong.Client.GetIp() returns the remote address of the client making the request.
// This will always return the address of the client directly connecting to Kong.
// That is, in cases when a load balancer is in front of Kong, this function will
// return the load balancer’s address, and not that of the downstream client.
func (c Client) GetIp() (ip string, err error) {
	ip_v, err := c.Ask(`kong.client.get_ip`)
	var ok bool
	if ip, ok = ip_v.(string); !ok {
		err = bridge.ReturnTypeError("string")
	}
	return
}

// kong.Client.GetForwardedIp() returns the remote address of the client making the request.
// Unlike kong.client.get_ip, this function will consider forwarded addresses in cases
// when a load balancer is in front of Kong. Whether this function returns a forwarded
// address or not depends on several Kong configuration parameters:
//
//   - trusted_ips
//   - real_ip_header
//   - real_ip_recursive
func (c Client) GetForwardedIp() (string, error) {
	return c.AskString(`kong.client.get_forwarded_ip`)
}

// kong.Client.GetPort() returns the remote port of the client making the request.
// This will always return the port of the client directly connecting to Kong.
// That is, in cases when a load balancer is in front of Kong, this function
// will return load balancer’s port, and not that of the downstream client.
func (c Client) GetPort() (string, error) {
	return c.AskString(`kong.client.get_port`)
}

// kong.Client.GetForwardedPort() returns the remote port of the client making the request.
// Unlike kong.client.get_port, this function will consider forwarded ports in cases
// when a load balancer is in front of Kong. Whether this function returns a forwarded
// port or not depends on several Kong configuration parameters:
//
//   - trusted_ips
//   - real_ip_header
//   - real_ip_recursive
func (c Client) GetForwardedPort() (string, error) {
	return c.AskString(`kong.client.get_forwarded_port`)
}

// kong.Client.GetCredential() returns the credentials of the currently authenticated consumer.
// If not set yet, it returns nil.
func (c Client) GetCredential() (cred AuthenticatedCredential, err error) {
	var val interface{}
	val, err = c.Ask(`kong.client.get_credential`)
	if err != nil {
		return
	}

	var ok bool
	fmt.Println(val)
	if cred, ok = val.(AuthenticatedCredential); !ok {
		err = bridge.ReturnTypeError("AuthenticatedCredential")
	}
	return
}

// kong.Client.LoadConsumer() returns the consumer from the datastore (or cache).
// Will look up the consumer by id, and optionally will do a second search by name.
func (c Client) LoadConsumer(consumer_id string, by_username bool) (consumer entities.Consumer, err error) {
	var reply interface{}
	reply, err = c.Ask(`kong.client.load_consumer`, consumer_id, by_username)
	if err != nil {
		return
	}

	return checkConsumer(reply)
}

// kong.Client.GetConsumer() returns the consumer entity of the currently authenticated consumer.
// If not set yet, it returns nil.
func (c Client) GetConsumer() (consumer entities.Consumer, err error) {
	var reply interface{}
	reply, err = c.Ask(`kong.client.get_consumer`)
	if err != nil {
		return
	}

	return checkConsumer(reply)
}

// kong.Client.Authenticate() sets the authenticated consumer and/or credential
// for the current request. While both consumer and credential can be nil,
// it is required that at least one of them exists. Otherwise this function will throw an error.
func (c Client) Authenticate(consumer *entities.Consumer, credential *AuthenticatedCredential) error {
	_, err := c.Ask(`kong.client.authenticate`, consumer, credential)
	return err
}

// kong.Client.GetProtocol() returns the protocol matched by the current route
// ("http", "https", "tcp" or "tls"), or nil, if no route has been matched,
// which can happen when dealing with erroneous requests.
func (c Client) GetProtocol(allow_terminated bool) (string, error) {
	return c.AskString(`kong.client.get_protocol`, allow_terminated)
}
