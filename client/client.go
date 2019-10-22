package client

import (
	"encoding/json"
	"fmt"
	"github.com/kong/go-pdk/entities"
)

type AuthenticatedCredential struct {
	Id         string `json:"id"`
	ConsumerId string `json:"consumer_id"`
}

type Client struct {
	ch chan string
}

func NewClient(ch chan string) *Client {
	return &Client{ch: ch}
}

func (c *Client) GetIp() string {
	c.ch <- `kong.client.get_ip`
	return <-c.ch
}

func (c *Client) GetForwardedIp() string {
	c.ch <- `kong.client.get_forwarded_ip`
	return <-c.ch
}

func (c *Client) GetPort() string {
	c.ch <- `kong.client.get_port`
	return <-c.ch
}

func (c *Client) GetForwardedPort() string {
	c.ch <- `kong.client.get_forwarded_port`
	return <-c.ch
}

func (c *Client) GetCredential() *AuthenticatedCredential {
	c.ch <- `kong.client.get_credential`
	reply := <-c.ch
	if reply == "null" {
		return nil
	}
	cred := AuthenticatedCredential{}
	json.Unmarshal([]byte(reply), &cred)
	return &cred
}

func (c *Client) LoadConsumer(consumer_id string, by_username bool) *entities.Consumer {
	c.ch <- `kong.client.load_consumer`
	reply := <-c.ch
	if reply == "null" {
		return nil
	}
	consumer := entities.Consumer{}
	json.Unmarshal([]byte(reply), &consumer)
	return &consumer
}

func (c *Client) GetConsumer() *entities.Consumer {
	c.ch <- `kong.client.get_consumer`
	reply := <-c.ch
	if reply == "null" {
		return nil
	}
	consumer := entities.Consumer{}
	json.Unmarshal([]byte(reply), &consumer)
	return &consumer
}

func (c *Client) Authenticate(consumer *entities.Consumer, credential *AuthenticatedCredential) error {
	consumerBytes, err := json.Marshal(consumer)
	if err != nil {
		return err
	}
	credBytes, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	c.ch <- fmt.Sprintf(`kong.client.authenticate:["%s","%s"]`,
		string(consumerBytes), string(credBytes))
	_ = <-c.ch

	return nil
}

func (c *Client) GetProtocol(allow_terminated bool) string {
	c.ch <- fmt.Sprintf(`kong.client.get_protocol:%t`, allow_terminated)
	return <-c.ch
}
