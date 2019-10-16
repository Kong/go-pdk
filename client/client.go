package client

import (
	"encoding/json"
	"github.com/kong/go-pdk/entities"
	"fmt"
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

func (c *Client) GetIp() *string {
	c.ch <- `kong.client.get_ip`
	ip := <-c.ch
	if ip == "null" {
		return nil
	}
	return &ip
}

func (c *Client) GetForwardedIp() *string {
	c.ch <- `kong.client.get_forwarded_ip`
	forwarded_ip := <-c.ch
	if forwarded_ip == "null" {
		return nil
	}
	return &forwarded_ip
}

func (c *Client) GetPort() *string {
	c.ch <- `kong.client.get_port`
	port := <-c.ch
	if port == "null" {
		return nil
	}
	return &port
}

func (c *Client) GetForwardedPort() *string {
	c.ch <- `kong.client.get_forwarded_port`
	forwarded_port := <-c.ch
	if forwarded_port == "null" {
		return nil
	}
	return &forwarded_port
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

// TODO client.authenticate

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

func (c *Client) GetProtocol(allow_terminated bool) *string {
	c.ch <- fmt.Sprintf(`kong.client.get_protocol:%b`, allow_terminated)
	reply := <-c.ch
	if reply == "null" {
		return nil
	}

	return &reply
}
