package request

import (
	"github.com/Kong/go-pdk/bridge"
)

type Request struct {
	bridge.PdkBridge
}

func New(ch chan string) Request {
	return Request{bridge.New(ch)}
}

func (r Request) GetScheme() (string, error) {
	return r.Ask(`kong.request.get_scheme`)
}

func (r Request) GetHost() (string, error) {
	return r.Ask(`kong.request.get_host`)
}

func (r Request) GetPort() (string, error) {
	return r.Ask(`kong.request.get_port`)
}

func (r Request) GetForwardedScheme() (string, error) {
	return r.Ask(`kong.request.get_forwarded_scheme`)
}

func (r Request) GetForwardedHost() (string, error) {
	return r.Ask(`kong.request.get_forwarded_host`)
}

func (r Request) GetForwardedPort() (string, error) {
	return r.Ask(`kong.request.get_forwarded_port`)
}

func (r Request) GetHttpVersion() (string, error) {
	return r.Ask(`kong.request.get_http_version`)
}

func (r Request) GetMethod() (string, error) {
	return r.Ask(`kong.request.get_method`)
}

func (r Request) GetPath() (string, error) {
	return r.Ask(`kong.request.get_path`)
}

func (r Request) GetPathWithQuery() (string, error) {
	return r.Ask(`kong.request.get_path_with_query`)
}

func (r Request) GetRawQuery() (string, error) {
	return r.Ask(`kong.request.get_raw_query`)
}

func (r Request) GetQueryArg() (string, error) {
	return r.Ask(`kong.request.get_query_arg`)
}

func (r Request) GetQuery(max_args int) (map[string]interface{}, error) {
	var res string
	var err error
	if max_args == -1 {
		res, err = r.Ask("kong.request.get_query")
	} else {
		res, err = r.Ask("kong.request.get_query", max_args)
	}
	if err != nil {
		return nil, err
	}

	query := make(map[string]interface{})
	bridge.Unmarshal(res, &query)
	return query, nil
}

func (r Request) GetHeader(k string) (string, error) {
	return r.Ask(`kong.request.get_header`, k)
}

func (r Request) GetHeaders(max_headers int) (map[string]interface{}, error) {
	var res string
	var err error
	if max_headers == -1 {
		res, err = r.Ask(`kong.request.get_headers`)
	} else {
		res, err = r.Ask(`kong.request.get_headers`, max_headers)
	}
	if err != nil {
		return nil, err
	}

	headers := make(map[string]interface{})
	bridge.Unmarshal(res, &headers)
	return headers, nil
}

func (r Request) GetRawBody() (string, error) {
	return r.Ask(`kong.request.get_raw_body`)
}

// TODO get_body
