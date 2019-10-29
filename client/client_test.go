package client

import (
	"testing"

	"github.com/kong/go-pdk/entities"
	"github.com/stretchr/testify/assert"
)

var client Client
var ch chan string

func init() {
	ch = make(chan string)
	client = New(ch)
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
	assert.Equal(t, "kong.client.get_ip:null", getName(func() { client.GetIp() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := client.GetIp(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetIp(); res <- r }, ""))
}

func TestGetForwardedIp(t *testing.T) {
	assert.Equal(t, "kong.client.get_forwarded_ip:null", getName(func() { client.GetForwardedIp() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := client.GetForwardedIp(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetForwardedIp(); res <- r }, ""))
}

func TestGetPort(t *testing.T) {
	assert.Equal(t, "kong.client.get_port:null", getName(func() { client.GetPort() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := client.GetPort(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetPort(); res <- r }, ""))
}

func TestGetForwardedPort(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetForwardedPort() }), "kong.client.get_forwarded_port:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := client.GetForwardedPort(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := client.GetForwardedPort(); res <- r }, ""), "")
}

func TestGetCredential(t *testing.T) {
	assert.Equal(t, "kong.client.get_credential:null", getName(func() { client.GetCredential() }))

	res := make(chan *AuthenticatedCredential)
	go func(res chan *AuthenticatedCredential) { r, _ := client.GetCredential(); res <- r }(res)
	_ = <-ch
	ch <- `{"id": "123", "consumer_id": "321"}`
	cred := <-res
	assert.Equal(t, "123", cred.Id)
	assert.Equal(t, "321", cred.ConsumerId)
}

func TestLoadConsumer(t *testing.T) {
	assert.Equal(t, "kong.client.get_credential:null", getName(func() { client.GetCredential() }))

	res := make(chan *entities.Consumer)
	go func(res chan *entities.Consumer) { r, _ := client.LoadConsumer("foo", true); res <- r }(res)
	_ = <-ch
	ch <- `
		{
			"id": "foo_id",
			"created_at": 123456,
			"username": "consumer",
			"custom_id": "custom_id"
		}`
	consumer := <-res
	assert.Equal(t, "foo_id", consumer.Id)
	assert.Equal(t, 123456, consumer.CreatedAt)
	assert.Equal(t, "consumer", *consumer.Username)
	assert.Equal(t, "custom_id", *consumer.CustomId)
}

func TestGetConsumer(t *testing.T) {
	assert.Equal(t, getName(func() { client.GetConsumer() }), "kong.client.get_consumer:null")

	res := make(chan *entities.Consumer)
	go func(res chan *entities.Consumer) { r, _ := client.GetConsumer(); res <- r }(res)
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

	expected := `kong.client.authenticate:[{"id":"","created_at":0,"username":"gs","custom_id":null,"tags":null},{"id":"gs","consumer_id":""}]`

	assert.Equal(t, expected, getName(func() {
		client.Authenticate(
			&entities.Consumer{Username: &n},
			&AuthenticatedCredential{Id: n})
	}))
}

func TestGetProtocol(t *testing.T) {
	assert.Equal(t, "kong.client.get_protocol:[true]", getName(func() { client.GetProtocol(true) }))
	assert.Equal(t, "kong.client.get_protocol:[false]", getName(func() { client.GetProtocol(false) }))
	assert.Equal(t, "tls", getStrValue(func(res chan string) { r, _ := client.GetForwardedPort(); res <- r }, "tls"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetForwardedPort(); res <- r }, ""))
}
