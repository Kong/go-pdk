package client

import (
	"encoding/json"
	"github.com/kong/go-pdk/entities"
)

type AuthenticatedCredential struct {
	Id         string
	ConsumerId string
}

type Client struct {
	ch chan string
}

func NewClient(ch chan string) *Client {
	return &Client{ch: ch}
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
