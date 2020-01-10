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

func (r Request) SetScheme(scheme string) error {
	_, err := r.Ask(`kong.service.request.set_scheme`, scheme)
	return err
}

func (r Request) SetPath(path string) error {
	_, err := r.Ask(`kong.service.request.set_path`, path)
	return err
}

// kong.ServiceRequest.SetRawQuery() sets the querystring
// of the request to the Service. The query argument is a string
// (without the leading ? character), and will not be processed in any way.
//
// For a higher-level function to set the query string from a ???? of arguments,
// see kong.ServiceRequest.SetQuery().
func (r Request) SetRawQuery(query string) error {
	_, err := r.Ask(`kong.service.request.set_raw_query`, query)
	return err
}

func (r Request) SetMethod(method string) error {
	_, err := r.Ask(`kong.service.request.set_method`, method)
	return err
}

// kong.ServiceRequest.SetQuery() sets the querystring of the request to the Service.
//
// Unlike kong.ServiceRequest.SetRawQuery(), the query argument must be a map
// in which each key is a string (corresponding to an arguments name), and each
// value is either a boolean, a string or an array of strings or booleans.
// Additionally, all string values will be URL-encoded.
//
// The resulting querystring will contain keys in their lexicographical order.
// The order of entries within the same key (when values are given as an array) is retained.
//
// If further control of the querystring generation is needed, a raw querystring
// can be given as a string with kong.ServiceRequest.SetRawQuery().
func (r Request) SetQuery(query map[string]interface{}) error {
	_, err := r.Ask(`kong.service.request.set_query`, query)
	return err
}

func (r Request) SetHeader(name string, value string) error {
	_, err := r.Ask(`kong.service.request.set_header`, name, value)
	return err
}

func (r Request) AddHeader(name string, value string) error {
	_, err := r.Ask(`kong.service.request.add_header`, name, value)
	return err
}

func (r Request) ClearHeader(name string) error {
	_, err := r.Ask(`kong.service.request.clear_header`, name)
	return err
}

func (r Request) SetHeaders(headers map[string]interface{}) error {
	_, err := r.Ask(`kong.service.request.set_headers`, headers)
	return err
}

// kong.ServiceRequest SetRawBody() sets the body of the request to the Service.
//
// The body argument must be a string and will not be processed in any way.
// This function also sets the Content-Length header appropriately.
// To set an empty body, one can give an empty string "" to this function.
//
// For a higher-level function to set the body based on the request content type,
// see kong.ServiceRequest.SetBody().
func (r Request) SetRawBody(body string) error {
	_, err := r.Ask(`kong.service.request.set_raw_body`, body)
	return err
}

// TODO set_body
