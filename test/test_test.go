package test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type foo struct {
	shouldExit bool
}

func MockNew() *foo {
	return &foo{
		// manually change me to see both ways
		shouldExit: true,
	}
}

func (f *foo) Response(kong *pdk.PDK) {
	if f.shouldExit {
		kong.Response.Exit(http.StatusInternalServerError, []byte(errors.New("exit").Error()), nil)
	}
}

func TestNoHangingChannel(t *testing.T) {
	t.Parallel()

	env, err := New(t, Request{
		Method: "POST",
		Url:    "",
		Body:   []byte("{}"),
	})
	assert.NoError(t, err)
	env.DoHttps(MockNew())
}

func TestSharedContext(t *testing.T) {
	t.Parallel()

	env, err := New(t, Request{
		Method: "POST",
		Url:    "",
		Body:   []byte("{}"),
	})
	assert.NoError(t, err)

	// Test set
	value, err := structpb.NewValue("test-token")
	assert.NoError(t, err)
	setArgs, err := proto.Marshal(&kong_plugin_protocol.KV{K: "Token", V: value})
	assert.NoError(t, err)

	env.Handle("kong.ctx.shared.set", setArgs)
	// Assert kv pair is in Ctx.Store
	assert.Equal(t, "test-token", env.Ctx.Store["Token"])

	// Test get
	key := bridge.WrapString("Token")
	getArgs, _ := proto.Marshal(key)
	response := env.Handle("kong.ctx.shared.get", getArgs)

	out := new(structpb.Value)
	err = proto.Unmarshal(response, out)
	assert.NoError(t, err)
	assert.Equal(t, "test-token", out.AsInterface())
}
