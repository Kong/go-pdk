package request

import (
	"github.com/Kong/go-pdk/bridge"
)

type Request struct {
	bridge.PdkBridge
}

func New(ch chan interface{}) Request {
	return Request{bridge.New(ch)}
}

func (r Request) GetScheme() (s string, err error) {
	return r.AskString(`kong.request.get_scheme`)
}

func (r Request) GetHost() (host string, err error) {
	return r.AskString(`kong.request.get_host`)
}

func (r Request) GetPort() (port string, err error) {
	return r.AskString(`kong.request.get_port`)
}

func (r Request) GetForwardedScheme() (s string, err error) {
	return r.AskString(`kong.request.get_forwarded_scheme`)
}

func (r Request) GetForwardedHost() (host string, err error) {
	return r.AskString(`kong.request.get_forwarded_host`)
}

func (r Request) GetForwardedPort() (port string, err error) {
	return r.AskString(`kong.request.get_forwarded_port`)
}

func (r Request) GetHttpVersion() (version string, err error) {
	return r.AskString(`kong.request.get_http_version`)
}

func (r Request) GetMethod() (m string, err error) {
	return r.AskString(`kong.request.get_method`)
}

func (r Request) GetPath() (string, error) {
	return r.AskString(`kong.request.get_path`)
}

func (r Request) GetPathWithQuery() (string, error) {
	return r.AskString(`kong.request.get_path_with_query`)
}

func (r Request) GetRawQuery() (string, error) {
	return r.AskString(`kong.request.get_raw_query`)
}

func (r Request) GetQueryArg() (string, error) {
	return r.AskString(`kong.request.get_query_arg`)
}

func (r Request) GetQuery(max_args int) (map[string]interface{}, error) {
	if max_args == -1 {
		return r.AskMap("kong.request.get_query")
	}

	return r.AskMap("kong.request.get_query", max_args)
}

func (r Request) GetHeader(k string) (string, error) {
	return r.AskString(`kong.request.get_header`, k)
}

func (r Request) GetHeaders(max_headers int) (map[string]interface{}, error) {
	if max_headers == -1 {
		return r.AskMap(`kong.request.get_headers`)
	}

	return r.AskMap(`kong.request.get_headers`, max_headers)
}

func (r Request) GetRawBody() (string, error) {
	return r.AskString(`kong.request.get_raw_body`)
}

// TODO get_body
