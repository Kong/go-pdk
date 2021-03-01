/*
Client request module.

A set of functions to retrieve information about the incoming requests made by clients.
*/
package request

import (
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"google.golang.org/protobuf/types/known/structpb"
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.Request`
type Request struct {
	bridge.PdkBridge
}

// Called by the plugin server at initialization.
// func New(ch chan interface{}) Request {
// 	return Request{bridge.New(ch)}
// }

// kong.Request.GetScheme() returns the scheme component of the request’s URL.
// The returned value is normalized to lower-case form.
func (r Request) GetScheme() (s string, err error) {
	return r.AskString(`kong.request.get_scheme`, nil)
}

// kong.Request.GetHost() returns the host component of the request’s URL,
// or the value of the “Host” header. The returned value is normalized
// to lower-case form.
func (r Request) GetHost() (host string, err error) {
	return r.AskString(`kong.request.get_host`, nil)
}

// kong.Request.GetPort() returns the port component of the request’s URL.
func (r Request) GetPort() (int, error) {
	return r.AskInt(`kong.request.get_port`, nil)
}

// kong.Request.GetForwardedScheme() returns the scheme component
// of the request’s URL, but also considers X-Forwarded-Proto if it
// comes from a trusted source. The returned value is normalized to lower-case.
//
// Whether this function considers X-Forwarded-Proto or not depends
// on several Kong configuration parameters:
//
//   - trusted_ips
//   - real_ip_header
//   - real_ip_recursive
//
// Note: support for the Forwarded HTTP Extension (RFC 7239) is not offered yet
// since it is not supported by ngx_http_realip_module.
func (r Request) GetForwardedScheme() (s string, err error) {
	return r.AskString(`kong.request.get_forwarded_scheme`, nil)
}

// kong.Request.GetForwardedHost() returns the host component of the request’s URL
// or the value of the “host” header. Unlike kong.Request.GetHost(), this function
// will also consider X-Forwarded-Host if it comes from a trusted source.
// The returned value is normalized to lower-case.
//
// Whether this function considers X-Forwarded-Proto or not depends
// on several Kong configuration parameters:
//
//   - trusted_ips
//   - real_ip_header
//   - real_ip_recursive
//
// Note: we do not currently offer support for Forwarded HTTP Extension (RFC 7239)
// since it is not supported by ngx_http_realip_module.
func (r Request) GetForwardedHost() (host string, err error) {
	return r.AskString(`kong.request.get_forwarded_host`, nil)
}

// kong.Request.GetForwardedPort() returns the port component of the request’s URL,
// but also considers X-Forwarded-Host if it comes from a trusted source.
//
// Whether this function considers X-Forwarded-Proto or not depends
// on several Kong configuration parameters:
//
//   - trusted_ips
//   - real_ip_header
//   - real_ip_recursive
//
// Note: we do not currently offer support for Forwarded HTTP Extension (RFC 7239)
// since it is not supported by ngx_http_realip_module.
func (r Request) GetForwardedPort() (int, error) {
	return r.AskInt(`kong.request.get_forwarded_port`, nil)
}

// kong.Request.GetHttpVersion() returns the HTTP version
// used by the client in the request, returning values
// such as "1"", "1.1", "2.0", or nil for unrecognized values.
func (r Request) GetHttpVersion() (version float64, err error) {
	return r.AskNumber(`kong.request.get_http_version`, nil)
}

// kong.Request.GetMethod() returns the HTTP method of the request.
// The value is normalized to upper-case.
func (r Request) GetMethod() (m string, err error) {
	return r.AskString(`kong.request.get_method`, nil)
}

// kong.Request.GetPath() returns the path component of the request’s URL.
// It is not normalized in any way and does not include the querystring.
func (r Request) GetPath() (string, error) {
	return r.AskString(`kong.request.get_path`, nil)
}

// kong.Request.GetPathWithQuery() returns the path, including
// the querystring if any. No transformations/normalizations are done.
func (r Request) GetPathWithQuery() (string, error) {
	return r.AskString(`kong.request.get_path_with_query`, nil)
}

// kong.Request.GetRawQuery() returns the query component of the request’s URL.
// It is not normalized in any way (not even URL-decoding of special characters)
// and does not include the leading ? character.
func (r Request) GetRawQuery() (string, error) {
	return r.AskString(`kong.request.get_raw_query`, nil)
}

// kong.Request.GetQueryArg() returns the value of the specified argument,
// obtained from the query arguments of the current request.
//
// The returned value is either a string, a boolean true if
// an argument was not given a value, or nil if no argument with name was found.
//
// If an argument with the same name is present multiple times in the querystring,
// this function will return the value of the first occurrence.
func (r Request) GetQueryArg(k string) (string, error) {
	return r.AskString(`kong.request.get_query_arg`, bridge.WrapString(k))
}

// kong.Request.GetQuery() returns a map of query arguments
// obtained from the querystring. Keys are query argument names.
// Values are either a string with the argument value, a boolean true
// if an argument was not given a value, or an array if an argument
// was given in the query string multiple times. Keys and values are
// unescaped according to URL-encoded escaping rules.
//
// Note that a query string `?foo&bar` translates to two boolean true arguments,
// and ?foo=&bar= translates to two string arguments containing empty strings.
//
// The max_args argument specifies the maximum number of returned arguments.
// Must be greater than 1 and not greater than 1000, or -1 to specify the
// default limit of 100 arguments.
func (r Request) GetQuery(max_args int) (map[string][]string, error) {
	if max_args == -1 {
		max_args = 100
	}

	arg := kong_plugin_protocol.Int{ V: int32(max_args) }
	out := new(structpb.Struct)
	err := r.Ask("kong.request.get_query", &arg, out)
	if err != nil {
		return nil, err
	}

	return bridge.UnwrapHeaders(out), nil
}

// kong.Request.GetHeader() returns the value of the specified request header.
//
// The returned value is either a string, or can be nil if a header with name
// was not found in the request. If a header with the same name is present
// multiple times in the request, this function will return the value of the
// first occurrence of this header.
//
// Header names in are case-insensitive and are normalized to lowercase,
// and dashes (-) can be written as underscores (_); that is, the header
// X-Custom-Header can also be retrieved as x_custom_header.
func (r Request) GetHeader(k string) (string, error) {
	return r.AskString(`kong.request.get_header`, bridge.WrapString(k))
}

// kong.Request.GetHeaders() returns a map holding the request headers.
// Keys are header names. Values are either a string with the header value,
// or an array of strings if a header was sent multiple times. Header names
// in this table are case-insensitive and are normalized to lowercase,
// and dashes (-) can be written as underscores (_); that is, the header
// X-Custom-Header can also be retrieved as x_custom_header.
//
// The max_args argument specifies the maximum number of returned headers.
// Must be greater than 1 and not greater than 1000, or -1 to specify the
// default limit of 100 headers.
func (r Request) GetHeaders(max_headers int) (map[string][]string, error) {
	if max_headers == -1 {
		max_headers = 100
	}

	arg := kong_plugin_protocol.Int{ V: int32(max_headers) }
	out := new(structpb.Struct)
	err := r.Ask("kong.request.get_headers", &arg, out)
	if err != nil {
		return nil, err
	}

	return bridge.UnwrapHeaders(out), nil
}

// kong.Request.GetRawBody() returns the plain request body.
//
// If the body has no size (empty), this function returns an empty string.
//
// If the size of the body is greater than the Nginx buffer size
// (set by client_body_buffer_size), this function will fail
// and return an error message explaining this limitation.
func (r Request) GetRawBody() (string, error) {
	return r.AskString(`kong.request.get_raw_body`, nil)
}

// TODO get_body
