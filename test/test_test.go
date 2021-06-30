package test

import (
	"net/http"
	"testing"

	"github.com/Kong/go-pdk"
	"github.com/stretchr/testify/assert"
)

type Config struct{}

// Plugin used for tests
func (conf *Config) Access(kong *pdk.PDK) {
	path, err := kong.Request.GetPath()
	if err != nil {
		kong.Response.Exit(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	switch path {
	case "/method":
		method, err := kong.Request.GetMethod()
		if err != nil {
			kong.Response.Exit(http.StatusInternalServerError, err.Error(), nil)
			return
		}

		switch method {
		case "GET":
			kong.Response.Exit(http.StatusOK, "get", nil)
		case "POST":
			kong.Response.Exit(http.StatusOK, "post", nil)
		default:
			kong.Response.ExitStatus(http.StatusNotImplemented)
		}
	default:
		kong.Response.ExitStatus(http.StatusNotImplemented)
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
	assert.Equal(t, "get", env.ClientRes.Body)
}

func TestAllowPOST(t *testing.T) {
	env, err := New(t, Request{
		Method: "POST",
		Url:    "http://localhost/method",
	})
	assert.NoError(t, err)

	env.DoHttps(&Config{})
	assert.Equal(t, 200, env.ClientRes.Status)
	assert.Equal(t, "post", env.ClientRes.Body)
}

func TestExitStatus(t *testing.T) {
	env, err := New(t, Request{
		Method: "POST",
		Url:    "http://localhost/notimplimented",
		Body:   "Should not copy",
	})
	assert.NoError(t, err)

	env.DoHttps(&Config{})
	assert.Equal(t, http.StatusNotImplemented, env.ClientRes.Status)
	assert.Equal(t, "", env.ClientRes.Body)
}
