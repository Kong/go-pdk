package nginx

import (
	"github.com/Kong/go-pdk/bridge"
)

func checkFloat(v interface{}) (f float64, err error) {
	f, ok := v.(float64)
	if !ok {
		err = bridge.ReturnTypeError("float64")
	}
	return
}

type Nginx struct {
	bridge.PdkBridge
}

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

func (n Nginx) GetCtxFloat(k string) (f float64, err error) {
	val, err := n.Ask(`kong.nginx.get_ctx`, k)
	if err != nil {
		return
	}

	return checkFloat(val)
}

func (n Nginx) ReqStartTime() (f float64, err error) {
	val, err := n.Ask(`kong.nginx.req_start_time`)
	if err != nil {
		return
	}

	return checkFloat(val)
}
