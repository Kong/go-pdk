package test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Kong/go-pdk"
	"github.com/stretchr/testify/assert"
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
