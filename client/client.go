package client

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
)

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

func New(ch chan interface{}) Client {
	return Client{bridge.New(ch)}
}

func (c Client) GetIp() (ip string, err error) {
	ip_v, err := c.Ask(`kong.client.get_ip`)
	var ok bool
	if ip, ok = ip_v.(string); !ok {
		err = bridge.ReturnTypeError("string")
	}
	return
}

func (c Client) GetForwardedIp() (string, error) {
	return c.AskString(`kong.client.get_forwarded_ip`)
}

func (c Client) GetPort() (string, error) {
	return c.AskString(`kong.client.get_port`)
}

func (c Client) GetForwardedPort() (string, error) {
	return c.AskString(`kong.client.get_forwarded_port`)
}

func (c Client) GetCredential() (cred AuthenticatedCredential, err error) {
	var val interface{}
	val, err = c.Ask(`kong.client.get_credential`)
	if err != nil {
		return
	}

	var ok bool
	if cred, ok = val.(AuthenticatedCredential); !ok {
		err = bridge.ReturnTypeError("AuthenticatedCredential")
	}
	return
}

func (c Client) LoadConsumer(consumer_id string, by_username bool) (consumer entities.Consumer, err error) {
	var reply interface{}
	reply, err = c.Ask(`kong.client.load_consumer`)
	if err != nil {
		return
	}

	return checkConsumer(reply)
}

func (c Client) GetConsumer() (consumer entities.Consumer, err error) {
	var reply interface{}
	reply, err = c.Ask(`kong.client.get_consumer`)
	if err != nil {
		return
	}

	return checkConsumer(reply)
}

func (c Client) Authenticate(consumer *entities.Consumer, credential *AuthenticatedCredential) error {
	_, err := c.Ask(`kong.client.authenticate`, consumer, credential)
	return err
}

func (c Client) GetProtocol(allow_terminated bool) (string, error) {
	return c.AskString(`kong.client.get_protocol`, allow_terminated)
}
