package client

import (
	"github.com/kong/go-pdk/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

var client *Client
var ch chan string

func init() {
	ch = make(chan string)
	client = &Client{ch: ch}
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func getStrValue(f func(res chan string), val string) string {
	res := make(chan string)
	go f(res)
	_ = <-ch
	ch <- val
	return <-res
}

func TestGetIp(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetIp() }), "kong.client.get_ip")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetIp() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetIp() }, ""), "")
}

func TestGetForwardedIp(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetForwardedIp() }), "kong.client.get_forwarded_ip")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetForwardedIp() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetForwardedIp() }, ""), "")
}

func TestGetPort(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetPort() }), "kong.client.get_port")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetPort() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetPort() }, ""), "")
}

func TestGetForwardedPort(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetForwardedPort() }), "kong.client.get_forwarded_port")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetForwardedPort() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetForwardedPort() }, ""), "")
}

func TestGetCredential(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetCredential() }), "kong.client.get_credential")

	res := make(chan *AuthenticatedCredential)
	go func(res chan *AuthenticatedCredential) { res <- client.GetCredential() }(res)
	_ = <-ch
	ch <- `{"id": "123", "consumer_id": "321"}`
	cred := <-res
	assert.Equal(t, cred.Id, "123")
	assert.Equal(t, cred.ConsumerId, "321")
}

func TestLoadConsumer(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetCredential() }), "kong.client.get_credential")

	res := make(chan *entities.Consumer)
	go func(res chan *entities.Consumer) { res <- client.LoadConsumer("foo", true) }(res)
	_ = <-ch
	ch <- `
		{
			"id": "foo_id",
			"created_at": 123456,
			"username": "consumer",
			"custom_id": "custom_id"
		}`
	consumer := <-res
	assert.Equal(t, consumer.Id, "foo_id")
	assert.Equal(t, consumer.CreatedAt, 123456)
	assert.Equal(t, *consumer.Username, "consumer")
	assert.Equal(t, *consumer.CustomId, "custom_id")
}

func TestGetConsumer(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetConsumer() }), "kong.client.get_consumer")

	res := make(chan *entities.Consumer)
	go func(res chan *entities.Consumer) { res <- client.GetConsumer() }(res)
	_ = <-ch
	ch <- `
		{
			"id": "foo_id",
			"created_at": 123456,
			"username": "consumer",
			"custom_id": "custom_id"
		}`
	consumer := <-res
	assert.Equal(t, consumer.Id, "foo_id")
	assert.Equal(t, consumer.CreatedAt, 123456)
	assert.Equal(t, *consumer.Username, "consumer")
	assert.Equal(t, *consumer.CustomId, "custom_id")

}

func TestAuthenticate(t *testing.T) {
	n := "gs"

	actual := `kong.client.authenticate:["{"id":"","created_at":0,"username":"gs","custom_id":null,"tags":null}","{"id":"gs","consumer_id":""}"]`

	assert.Equal(t, getName(func() {
		client.Authenticate(
			&entities.Consumer{Username: &n},
			&AuthenticatedCredential{Id: n})
	}), actual)
}

func TestGetProtocol(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetProtocol(true) }), "kong.client.get_protocol:true")
	assert.Equal(t, getName(func() { client.GetProtocol(false) }), "kong.client.get_protocol:false")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetForwardedPort() }, "tls"), "tls")
	assert.Equal(t, getStrValue(func(res chan string) { res <- client.GetForwardedPort() }, ""), "")
}
