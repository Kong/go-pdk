/*
Access Nginx APIs.
*/
package nginx

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"google.golang.org/protobuf/types/known/structpb"
)

// Holds this module's functions.  Accessible as `kong.Nginx`
type Nginx struct {
	bridge.PdkBridge
}

// kong.Nginx.GetVar() returns an Nginx variable.  Equivalent to `ngx.var[k]`
func (n Nginx) GetVar(k string) (string, error) {
	return n.AskString(`kong.nginx.get_var`, bridge.WrapString(k))
}

func (n Nginx) GetTLS1VersionStr() (string, error) {
	return n.AskString(`kong.nginx.get_tls1_version_str`, nil)
}

// kong.Nginx.SetCtx() sets a value in the `ngx.ctx` request context table.
func (n Nginx) SetCtx(k string, v interface{}) error {
	v2, err := structpb.NewValue(v)
	if err != nil {
		return err
	}
	arg := kong_plugin_protocol.KV{
		K: k,
		V: v2,
	}
	return n.Ask(`kong.nginx.set_ctx`, &arg, nil)
}

// kong.Nginx.GetCtxAny() returns a value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxAny(k string) (interface{}, error) {
	return n.AskValue(`kong.nginx.get_ctx`, bridge.WrapString(k))
}

// kong.Nginx.GetCtxString() returns a string value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxString(k string) (string, error) {
	out := new(structpb.Value)
	if err := n.Ask(`kong.nginx.get_ctx`, bridge.WrapString(k), out); err != nil {
		return "", err
	}
	return out.GetStringValue(), nil
}

// kong.Nginx.GetCtxFloat() returns a float value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxFloat(k string) (float64, error) {
	out := new(structpb.Value)
	if err := n.Ask(`kong.nginx.get_ctx`, bridge.WrapString(k), out); err != nil {
		return 0, err
	}
	return out.GetNumberValue(), nil
}

// kong.Nginx.GetCtxInt() returns an integer value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxInt(k string) (int, error) {
	out := new(structpb.Value)
	if err := n.Ask(`kong.nginx.get_ctx`, bridge.WrapString(k), out); err != nil {
		return 0, err
	}
	return int(out.GetNumberValue()), nil
}

// kong.Nginx.ReqStartTime() returns the curent request's start time
// as a floating-point number of seconds.  Equivalent to `ngx.req.start_time()`
func (n Nginx) ReqStartTime() (float64, error) {
	return n.AskNumber(`kong.nginx.req_start_time`, nil)
}

// kong.Nginx.GetSubsystem() returns the current Nginx subsystem
// this function is called from: “http” or “stream”.
func (n Nginx) GetSubsystem() (string, error) {
	return n.AskString(`kong.nginx.get_subsystem`, nil)
}
