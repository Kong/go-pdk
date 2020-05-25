/*
Access Nginx APIs.
*/
package nginx

import (
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.Nginx`
type Nginx struct {
	bridge.PdkBridge
}

// Called by the plugin server at initialization.
func New(ch chan interface{}) Nginx {
	return Nginx{bridge.New(ch)}
}

// kong.Nginx.GetVar() returns an Nginx variable.  Equivalent to `ngx.var[k]`
func (n Nginx) GetVar(k string) (string, error) {
	return n.AskString(`kong.nginx.get_var`, k)
}

func (n Nginx) GetTLS1VersionStr() (string, error) {
	return n.AskString(`kong.nginx.get_tls1_version_str`)
}

// kong.Nginx.SetCtx() sets a value in the `ngx.ctx` request context table.
func (n Nginx) SetCtx(k string, v interface{}) error {
	_, err := n.Ask(`kong.nginx.set_ctx`, k, v)
	return err
}

// kong.Nginx.GetCtxAny() returns a value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxAny(k string) (interface{}, error) {
	return n.Ask(`kong.nginx.get_ctx`, k)
}

// kong.Nginx.GetCtxString() returns a string value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxString(k string) (string, error) {
	return n.AskString(`kong.nginx.get_ctx`, k)
}

// kong.Nginx.GetCtxFloat() returns a float value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxFloat(k string) (float64, error) {
	return n.AskFloat(`kong.nginx.get_ctx`, k)
}

// kong.Nginx.GetCtxInt() returns an integer value from the `ngx.ctx` request context table.
func (n Nginx) GetCtxInt(k string) (int, error) {
	return n.AskInt(`kong.nginx.get_ctx`, k)
}

// kong.Nginx.ReqStartTime() returns the curent request's start time
// as a floating-point number of seconds.  Equivalent to `ngx.req.start_time()`
func (n Nginx) ReqStartTime() (float64, error) {
	return n.AskFloat(`kong.nginx.req_start_time`)
}

// kong.Nginx.GetSubsystem() returns the current Nginx subsystem
// this function is called from: “http” or “stream”.
func (n Nginx) GetSubsystem() (string, error) {
	return n.AskString(`kong.nginx.get_subsystem`)
}
