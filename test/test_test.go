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

type Config struct {
	shouldExit bool
}

// Plugin used for tests
func (conf *Config) Access(kong *pdk.PDK) {
	if conf.shouldExit {
		kong.Response.Exit(http.StatusInternalServerError, []byte(errors.New("exit").Error()), nil)
	}

	path, err := kong.Request.GetPath()
	if err != nil {
		kong.Response.Exit(http.StatusInternalServerError, []byte(err.Error()), nil)
		return
	}

	switch path {
	case "/method":
		method, err := kong.Request.GetMethod()
		if err != nil {
			kong.Response.Exit(http.StatusInternalServerError, []byte(err.Error()), nil)
			return
		}

		switch method {
		case "GET":
			kong.Response.Exit(http.StatusOK, []byte("get"), nil)
		case "POST":
			kong.Response.Exit(http.StatusOK, []byte("post"), nil)
		default:
			kong.Response.Exit(http.StatusNotImplemented, nil, nil)
		}
	default:
		kong.Response.Exit(http.StatusNotImplemented, nil, nil)
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
	env.DoHttps(&Config{ shouldExit: true })
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

func TestAllowGET(t *testing.T) {
	env, err := New(t, Request{
		Method: "GET",
		Url:    "http://localhost/method",
	})
	assert.NoError(t, err)

	env.DoHttps(&Config{})
	assert.Equal(t, 200, env.ClientRes.Status)
	assert.Equal(t, []byte("get"), env.ClientRes.Body)
}

func TestAllowPOST(t *testing.T) {
	env, err := New(t, Request{
		Method: "POST",
		Url:    "http://localhost/method",
		Body:   []byte("post body"),
	})
	assert.NoError(t, err)

	env.DoHttps(&Config{})
	assert.Equal(t, 200, env.ClientRes.Status)
	assert.Equal(t, []byte("post"), env.ClientRes.Body)
}

func TestExitStatus(t *testing.T) {
	t.Parallel()
	env, err := New(t, Request{
		Method: "POST",
		Url:    "http://localhost/notimplimented",
		Body:   []byte("Should not copy"),
	})
	assert.NoError(t, err)

	env.DoHttps(&Config{})
	assert.Equal(t, http.StatusNotImplemented, env.ClientRes.Status)
	assert.Equal(t, []byte(nil), env.ClientRes.Body)
}
