/*
Manipulation of the request to the Service.
*/
package request

import (
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.ServiceRequest`
type Request struct {
	bridge.PdkBridge
}

// Called by the plugin server at initialization.
func New(ch chan interface{}) Request {
	return Request{bridge.New(ch)}
}

// kong.ServiceRequest.SetScheme() sets the protocol to use
// when proxying the request to the Service.
// Supported values are "http" or "https".
func (r Request) SetScheme(scheme string) error {
	_, err := r.Ask(`kong.service.request.set_scheme`, scheme)
	return err
}

// kong.ServiceRequest.SetPath() sets the path component
// for the request to the service. It is not normalized
// in any way and should not include the querystring.
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

// kong.ServiceRequest.SetMethod() sets the HTTP method
// for the request to the service.
//
// Supported method values are: "GET", "HEAD", "PUT", "POST",
// "DELETE", "OPTIONS", "MKCOL", "COPY", "MOVE", "PROPFIND",
// "PROPPATCH", "LOCK", "UNLOCK", "PATCH", "TRACE".
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
func (r Request) SetQuery(query map[string][]string) error {
	_, err := r.Ask(`kong.service.request.set_query`, query)
	return err
}

// kong.ServiceRequest.SetHeader() sets a header in the request
// to the Service with the given value. Any existing header with
// the same name will be overridden.
//
// If the header argument is "host" (case-insensitive), then this
// will also set the SNI of the request to the Service.
func (r Request) SetHeader(name string, value string) error {
	_, err := r.Ask(`kong.service.request.set_header`, name, value)
	return err
}

// kong.ServiceRequest.AddHeader() adds a request header with the given value
// to the request to the Service. Unlike kong.ServiceRequest.SetHeader(),
// this function will not remove any existing headers with the same name.
// Instead, several occurences of the header will be present in the request.
// The order in which headers are added is retained.
func (r Request) AddHeader(name string, value string) error {
	_, err := r.Ask(`kong.service.request.add_header`, name, value)
	return err
}

// kong.ServiceRequest.ClearHeader() removes all occurrences
// of the specified header in the request to the Service.
func (r Request) ClearHeader(name string) error {
	_, err := r.Ask(`kong.service.request.clear_header`, name)
	return err
}

// kong.ServiceRequest.SetHeaders() sets the headers of the request
// to the Service. Unlike kong.ServiceRequest.SetHeader(), the headers argument
// must be a map in which each key is a string (corresponding to a header’s name),
// and each value is a string, or an array of strings.
//
// The resulting headers are produced in lexicographical order.
// The order of entries with the same name (when values are given as an array) is retained.
//
// This function overrides any existing header bearing the same name as those
// specified in the headers argument. Other headers remain unchanged.
//
// If the "Host" header is set (case-insensitive), then this is will also set
// the SNI of the request to the Service.
func (r Request) SetHeaders(headers map[string][]string) error {
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
