package client

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
)

type Client struct {
	bridge.PdkBridge
}

type AuthenticatedCredential struct {
	Id         string `json:"id"`
	ConsumerId string `json:"consumer_id"`
}

func New(ch chan string) Client {
	return Client{bridge.New(ch)}
}

func (c Client) GetIp() (string, error) {
	return c.Ask(`kong.client.get_ip`)
}

func (c Client) GetForwardedIp() (string, error) {
	return c.Ask(`kong.client.get_forwarded_ip`)
}

func (c Client) GetPort() (string, error) {
	return c.Ask(`kong.client.get_port`)
}

func (c Client) GetForwardedPort() (string, error) {
	return c.Ask(`kong.client.get_forwarded_port`)
}

func (c Client) GetCredential() (*AuthenticatedCredential, error) {
	if reply, err := c.Ask(`kong.client.get_credential`); err != nil {
		return nil, err
	} else {
		cred := AuthenticatedCredential{}
		bridge.Unmarshal(reply, &cred)
		return &cred, nil
	}
}

func (c Client) LoadConsumer(consumer_id string, by_username bool) (*entities.Consumer, error) {
	if reply, err := c.Ask(`kong.client.load_consumer`); err != nil {
		return nil, err
	} else {
		consumer := entities.Consumer{}
		bridge.Unmarshal(reply, &consumer)
		return &consumer, nil
	}
}

func (c Client) GetConsumer() (*entities.Consumer, error) {
	if reply, err := c.Ask(`kong.client.get_consumer`); err != nil {
		return nil, err
	} else {
		consumer := entities.Consumer{}
		bridge.Unmarshal(reply, &consumer)
		return &consumer, nil
	}
}

func (c Client) Authenticate(consumer *entities.Consumer, credential *AuthenticatedCredential) error {
	_, err := c.Ask(`kong.client.authenticate`, consumer, credential)
	return err
}

func (c Client) GetProtocol(allow_terminated bool) (string, error) {
	return c.Ask(`kong.client.get_protocol`, allow_terminated)
}
