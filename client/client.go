package client

import (
	"fmt"
	"github.com/kong/go-pdk/bridge"
	"github.com/kong/go-pdk/entities"
)

type Client struct {
	bridge.PdkBridge
}

type AuthenticatedCredential struct {
	Id         string `json:"id"`
	ConsumerId string `json:"consumer_id"`
}

func New(ch chan string) *Client {
	return &Client{*bridge.New(ch)}
}

func (c *Client) GetIp() string {
	return c.Ask(`kong.client.get_ip`)
}

func (c *Client) GetForwardedIp() string {
	return c.Ask(`kong.client.get_forwarded_ip`)
}

func (c *Client) GetPort() string {
	return c.Ask(`kong.client.get_port`)
}

func (c *Client) GetForwardedPort() string {
	return c.Ask(`kong.client.get_forwarded_port`)
}

func (c *Client) GetCredential() *AuthenticatedCredential {
	reply := c.Ask(`kong.client.get_credential`)
	if reply == "null" {
		return nil
	}
	cred := AuthenticatedCredential{}
	bridge.Unmarshal(reply, &cred)
	return &cred
}

func (c *Client) LoadConsumer(consumer_id string, by_username bool) *entities.Consumer {
	reply := c.Ask(`kong.client.load_consumer`)
	if reply == "null" {
		return nil
	}
	consumer := entities.Consumer{}
	bridge.Unmarshal(reply, &consumer)
	return &consumer
}

func (c *Client) GetConsumer() *entities.Consumer {
	reply := c.Ask(`kong.client.get_consumer`)
	if reply == "null" {
		return nil
	}
	consumer := entities.Consumer{}
	bridge.Unmarshal(reply, &consumer)
	return &consumer
}

func (c *Client) Authenticate(consumer *entities.Consumer, credential *AuthenticatedCredential) error {
	consumerBytes, err := bridge.Marshal(consumer)
	if err != nil {
		return err
	}
	credBytes, err := bridge.Marshal(credential)
	if err != nil {
		return err
	}

	_ = c.Ask(fmt.Sprintf(`kong.client.authenticate:["%s","%s"]`,
						  consumerBytes, credBytes))

	return nil
}

func (c *Client) GetProtocol(allow_terminated bool) string {
	return c.Ask(fmt.Sprintf(`kong.client.get_protocol:%t`, allow_terminated))
}
