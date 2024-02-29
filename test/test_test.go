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
	key := "TestKey"

	t.Parallel()

	env, err := New(t, Request{
		Method: "POST",
		Url:    "",
		Body:   []byte("{}"),
	})
	assert.NoError(t, err)

	perform := func(v interface{}) {
		// Test set
		value, err := structpb.NewValue(v)
		assert.NoError(t, err)
		setArgs, err := proto.Marshal(&kong_plugin_protocol.KV{K: key, V: value})
		assert.NoError(t, err)

		env.Handle("kong.ctx.shared.set", setArgs)
		// Assert kv pair is in Ctx.Store
		assert.Equal(t, v, env.Ctx.Store[key])

		// Test get
		keyWrapped := bridge.WrapString(key)
		getArgs, _ := proto.Marshal(keyWrapped)
		response := env.Handle("kong.ctx.shared.get", getArgs)

		out := new(structpb.Value)
		err = proto.Unmarshal(response, out)
		assert.NoError(t, err)
		assert.Equal(t, v, out.AsInterface())
	}

	testValues := []interface{}{
		"TestValue",
		3.14,
		false,
		map[string]interface{}{
			"key1": "value1",
		},
	}

	for _, v := range testValues {
		perform(v)
	}
}
