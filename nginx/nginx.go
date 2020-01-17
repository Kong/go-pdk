/*
Access Nginx variables.
*/
package nginx

import (
	"github.com/Kong/go-pdk/bridge"
)

type Nginx struct {
	bridge.PdkBridge
}

// Called by the plugin server at initialization.
func New(ch chan interface{}) Nginx {
	return Nginx{bridge.New(ch)}
}

func (n Nginx) GetVar(k string) (string, error) {
	return n.AskString(`kong.nginx.get_var`, k)
}

func (n Nginx) GetTLS1VersionStr() (string, error) {
	return n.AskString(`kong.nginx.get_tls1_version_str`)
}

func (n Nginx) GetCtxAny(k string) (interface{}, error) {
	return n.Ask(`kong.nginx.get_ctx`, k)
}

func (n Nginx) GetCtxString(k string) (string, error) {
	return n.AskString(`kong.nginx.get_ctx`, k)
}

func (n Nginx) GetCtxFloat(k string) (float64, error) {
	return n.AskFloat(`kong.nginx.get_ctx`, k)
}

func (n Nginx) ReqStartTime() (float64, error) {
	return n.AskFloat(`kong.nginx.req_start_time`)
}
