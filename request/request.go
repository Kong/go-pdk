package request

import (
	"fmt"
	"github.com/kong/go-pdk/bridge"
)

type Request struct {
	bridge.PdkBridge
}

func New(ch chan string) *Request {
	return &Request{*bridge.New(ch)}
}

func (r *Request) GetScheme() string {
	return r.Ask(`kong.request.get_scheme`)
}

func (r *Request) GetHost() string {
	return r.Ask(`kong.request.get_host`)
}

func (r *Request) GetPort() string {
	return r.Ask(`kong.request.get_port`)
}

func (r *Request) GetForwardedScheme() string {
	return r.Ask(`kong.request.get_forwarded_scheme`)
}

func (r *Request) GetForwardedHost() string {
	return r.Ask(`kong.request.get_forwarded_host`)
}

func (r *Request) GetForwardedPort() string {
	return r.Ask(`kong.request.get_forwarded_port`)
}

func (r *Request) GetHttpVersion() string {
	return r.Ask(`kong.request.get_http_version`)
}

func (r *Request) GetMethod() string {
	return r.Ask(`kong.request.get_method`)
}

func (r *Request) GetPath() string {
	return r.Ask(`kong.request.get_path`)
}

func (r *Request) GetPathWithQuery() string {
	return r.Ask(`kong.request.get_path_with_query`)
}

func (r *Request) GetRawQuery() string {
	return r.Ask(`kong.request.get_raw_query`)
}

func (r *Request) GetQueryArg() string {
	return r.Ask(`kong.request.get_query_arg`)
}

func (r *Request) GetQuery(max_args int) map[string]interface{} {
	var method string
	if max_args == -1 {
		method = "kong.request.get_query"
	} else {
		method = fmt.Sprintf(`kong.request.get_query:%d`, max_args)
	}

	query := make(map[string]interface{})
	bridge.Unmarshal(r.Ask(method), &query)
	return query
}

func (r *Request) GetHeader(k string) string {
	return r.Ask(fmt.Sprintf(`kong.request.get_header:%s`, k))
}

func (r *Request) GetHeaders(max_headers int) map[string]interface{} {
	var method string
	if max_headers == -1 {
		method = `kong.request.get_headers`
	} else {
		method = fmt.Sprintf(`kong.request.get_headers:%d`, max_headers)
	}

	headers := make(map[string]interface{})
	bridge.Unmarshal(r.Ask(method), &headers)
	return headers
}

func (r *Request) GetRawBody() string {
	return r.Ask(`kong.request.get_raw_body`)
}

// TODO get_body
