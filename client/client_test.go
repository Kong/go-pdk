package client

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/entities"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)

func mockClient(t *testing.T, s []bridgetest.MockStep) Client {
	return Client{bridge.New(bridgetest.Mock(t, s))}
}

func TestGetIp(t *testing.T) {
	c := mockClient(t, []bridgetest.MockStep{
		{"kong.client.get_ip", nil, bridge.WrapString("10.10.10.1")},
	})

	resp, err := c.GetIp()
	assert.NoError(t, err)
	assert.Equal(t, resp, "10.10.10.1")
}

func TestGetForwardedIp(t *testing.T) {
	c := mockClient(t, []bridgetest.MockStep{
		{"kong.client.get_forwarded_ip", nil, bridge.WrapString("10.10.10.1")},
	})

	resp, err := c.GetForwardedIp()
	assert.NoError(t, err)
	assert.Equal(t, resp, "10.10.10.1")
}

func TestGetPort(t *testing.T) {
	c := mockClient(t, []bridgetest.MockStep{
		{"kong.client.get_port", nil, &kong_plugin_protocol.Int{V: 443}},
	})
	resp, err := c.GetPort()
	assert.NoError(t, err)
	assert.Equal(t, 443, resp)
}

func TestGetForwardedPort(t *testing.T) {
	c := mockClient(t, []bridgetest.MockStep{
		{"kong.client.get_forwarded_port", nil, &kong_plugin_protocol.Int{V: 80}},
	})
	resp, err := c.GetForwardedPort()
	assert.NoError(t, err)
	assert.Equal(t, 80, resp)
}

func TestGetCredential(t *testing.T) {
	c := mockClient(t, []bridgetest.MockStep{
		{"kong.client.get_credential", nil,
			&kong_plugin_protocol.AuthenticatedCredential{Id: "000:00", ConsumerId: "000:01"},
		},
	})

	resp, err := c.GetCredential()
	assert.NoError(t, err)
	assert.Equal(t, AuthenticatedCredential{Id: "000:00", ConsumerId: "000:01"}, resp)
}

func TestLoadConsumer(t *testing.T) {
	c := mockClient(t, []bridgetest.MockStep{
		{"kong.client.load_consumer",
			&kong_plugin_protocol.ConsumerSpec{Id: "001", ByUsername: false},
			&kong_plugin_protocol.Consumer{Id: "001", Username: "Jon Doe"},
		},
	})
	resp, err := c.LoadConsumer("001", false)
	assert.NoError(t, err)
	assert.Equal(t, entities.Consumer{Id: "001", Username: "Jon Doe"}, resp)
}

/*
func TestGetConsumer(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.client.get_consumer"}, getBack(func() { client.GetConsumer() }))
}

func TestAuthenticate(t *testing.T) {
	var consumer *entities.Consumer = nil
	var credential *AuthenticatedCredential = nil
	assert.Equal(t, bridge.StepData{Method: "kong.client.authenticate", Args: []interface{}{consumer, credential}}, getBack(func() { client.Authenticate(nil, nil) }))
}

func TestGetProtocol(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.client.get_protocol", Args: []interface{}{true}}, getBack(func() { client.GetProtocol(true) }))
	assert.Equal(t, bridge.StepData{Method: "kong.client.get_protocol", Args: []interface{}{false}}, getBack(func() { client.GetProtocol(false) }))
}
*/
