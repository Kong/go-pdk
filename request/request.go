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

// kong.Request.GetPort() returns the port component of the request’s URL.
func (r Request) GetPort() (int, error) {
	return r.AskInt(`kong.request.get_port`)
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
	return r.AskString(`kong.request.get_forwarded_scheme`)
}

func (r Request) GetForwardedHost() (host string, err error) {
	return r.AskString(`kong.request.get_forwarded_host`)
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
	return r.AskInt(`kong.request.get_forwarded_port`)
}

// kong.Request.GetHttpVersion() returns the HTTP version
// used by the client in the request, returning values
// such as "1"", "1.1", "2.0", or nil for unrecognized values.
func (r Request) GetHttpVersion() (version float64, err error) {
	return r.AskFloat(`kong.request.get_http_version`)
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

// kong.Request.GetQueryArg() returns the value of the specified argument,
// obtained from the query arguments of the current request.
//
// The returned value is either a string, a boolean true if
// an argument was not given a value, or nil if no argument with name was found.
//
// If an argument with the same name is present multiple times in the querystring,
// this function will return the value of the first occurrence.
func (r Request) GetQueryArg(k string) (string, error) {
	return r.AskString(`kong.request.get_query_arg`, k)
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
func (r Request) GetQuery(max_args int) (map[string]interface{}, error) {
	if max_args == -1 {
		return r.AskMap("kong.request.get_query")
	}

	return r.AskMap("kong.request.get_query", max_args)
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
	return r.AskString(`kong.request.get_header`, k)
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
