package bridge

import (
	"testing"

	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
)


func TestAsk(t *testing.T) {
	b := New(bridgetest.Mock(t, []bridgetest.MockStep{
		{"foo.bar", WrapString("first"), WrapString("resp")},
	}))


	out := new(kong_plugin_protocol.String)
	err := b.Ask("foo.bar", WrapString("first"), out)
	if err != nil {
		t.Fatalf("got this: %s", err)
	}
	if out.V != "resp" {
		t.Fatalf("no 'resp': %v", out.V)
	}
	b.Close()
}

func TestAskString(t *testing.T) {
	b := New(bridgetest.Mock(t, []bridgetest.MockStep{
		{"foo.bar", WrapString("first"), WrapString("resp")},
	}))

	ret, err := b.AskString("foo.bar", WrapString("first"))
	if err != nil {
		t.Fatalf("got this: %s", err)
	}
	if ret != "resp" {
		t.Fatalf("no 'resp': %v", ret)
	}
}
