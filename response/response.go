/*
Client response module.

The downstream response module contains a set of functions for producing
and manipulating responses sent back to the client (“downstream”).
Responses can be produced by Kong (e.g. an authentication plugin rejecting
a request), or proxied back from an Service’s response body.

Unlike kong.ServiceResponse, this module allows mutating the response
before sending it back to the client.
*/
package response

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"google.golang.org/protobuf/types/known/structpb"
)

// Holds this module's functions.  Accessible as `kong.Response`
type Response struct {
	bridge.PdkBridge
}

// kong.Response.GetStatus() returns the HTTP status code
// currently set for the downstream response (as an integer).
//
// If the request was proxied (as per kong.Response.GetSource()),
// the return value will be that of the response from the Service
// (identical to kong.ServiceResponse.GetStatus()).
//
// If the request was not proxied, and the response was produced
// by Kong itself (i.e. via kong.Response.Exit()), the return value
// will be returned as-is.
func (r Response) GetStatus() (int, error) {
	return r.AskInt(`kong.response.get_status`, nil)
}

// kong.Response.GetHeader() returns the value of the specified
// response header, as would be seen by the client once received.
//
// The list of headers returned by this function can consist of
// both response headers from the proxied Service and headers
// added by Kong (e.g. via kong.Response.AddHeader()).
//
// The return value is either a string, or can be nil if a header
// with name was not found in the response. If a header with the
// same name is present multiple times in the request, this function
// will return the value of the first occurrence of this header.
//
// Header names are case-insensitive and dashes (-) can be written
// as underscores (_); that is, the header X-Custom-Header
// can also be retrieved as x_custom_header.
func (r Response) GetHeader(name string) (string, error) {
	return r.AskString(`kong.response.get_header`, bridge.WrapString(name))
}

// kong.Response.GetHeaders() returns a map holding the response headers.
// Keys are header names. Values are either a string with the header value,
// or an array of strings if a header was sent multiple times.
// Header names in this table are case-insensitive and are normalized
// to lowercase, and dashes (-) can be written as underscores (_);
// that is, the header X-Custom-Header can also be retrieved as x_custom_header.
//
// A response initially has no headers until a plugin short-circuits
// the proxying by producing one (e.g. an authentication plugin rejecting
// a request), or the request has been proxied, and one of the latter
// execution phases is currently running.
//
// Unlike kong.ServiceResponse.GetHeaders(), this function returns
// all headers as the client would see them upon reception,
// including headers added by Kong itself.
//
// The max_args argument specifies the maximum number of returned headers.
// Must be greater than 1 and not greater than 1000, or -1 to specify the
// default limit of 100 arguments.
func (r Response) GetHeaders(max_headers int) (res map[string][]string, err error) {
	if max_headers == -1 {
		max_headers = 100
	}

	arg := kong_plugin_protocol.Int{V: int32(max_headers)}
	out := new(structpb.Struct)
	err = r.Ask(`kong.response.get_headers`, &arg, out)
	if err != nil {
		return nil, err
	}

	return bridge.UnwrapHeaders(out), nil
}

// kong.Response.GetSource() helps determining where the current response
// originated from. Kong being a reverse proxy, it can short-circuit
// a request and produce a response of its own, or the response can
// come from the proxied Service.
//
// Returns a string with three possible values:
//
// - “exit” is returned when, at some point during the processing of the request,
// there has been a call to kong.response.exit(). In other words, when the request
// was short-circuited by a plugin or by Kong itself (e.g. invalid credentials).
//
// - “error” is returned when an error has happened while processing the request
// - for example, a timeout while connecting to the upstream service.
//
// - “service” is returned when the response was originated by
// successfully contacting the proxied Service.
func (r Response) GetSource() (string, error) {
	return r.AskString(`kong.response.get_source`, nil)
}

// kong.Response.SetStatus() allows changing the downstream response
// HTTP status code before sending it to the client.
//
// This function should be used in the header_filter phase,
// as Kong is preparing headers to be sent back to the client.
func (r Response) SetStatus(status int) error {
	arg := kong_plugin_protocol.Int{V: int32(status)}
	err := r.Ask(`kong.response.set_status`, &arg, nil)
	return err
}

// kong.Response.SetHeader() sets a response header with the given value.
// This function overrides any existing header with the same name.
//
// This function should be used in the header_filter phase,
// as Kong is preparing headers to be sent back to the client.
func (r Response) SetHeader(k string, v string) error {
	arg := kong_plugin_protocol.KV{
		K: k,
		V: structpb.NewStringValue(v),
	}
	err := r.Ask(`kong.response.set_header`, &arg, nil)
	return err
}

// kong.Response.AddHeader() adds a response header with the given value.
// Unlike kong.Response.SetHeader(), this function does not remove
// any existing header with the same name. Instead, another header
// with the same name will be added to the response. If no header
// with this name already exists on the response, then it is added
// with the given value, similarly to kong.Response.SetHeader().
//
// This function should be used in the header_filter phase,
// as Kong is preparing headers to be sent back to the client.
func (r Response) AddHeader(k string, v string) error {
	arg := kong_plugin_protocol.KV{
		K: k,
		V: structpb.NewStringValue(v),
	}
	err := r.Ask(`kong.response.add_header`, &arg, nil)
	return err
}

// kong.Response.ClearHeader() removes all occurrences of the specified header
// in the response sent to the client.
//
// This function should be used in the header_filter phase,
// as Kong is preparing headers to be sent back to the client.
func (r Response) ClearHeader(k string) error {
	err := r.Ask(`kong.response.clear_header`, bridge.WrapString(k), nil)
	return err
}

// kong.Response.SetHeaders() sets the headers for the response.
// Unlike kong.Response.SetHeader(), the headers argument must be a map
// in which each key is a string (corresponding to a header’s name),
// and each value is an array of strings.  To clear a previously
// set header, you can set it's value to an empty array.
//
// This function should be used in the header_filter phase,
// as Kong is preparing headers to be sent back to the client.
//
// The resulting headers are produced in lexicographical order.
// The order of entries with the same name (when values are given
// as an array) is retained.
//
// This function overrides any existing header bearing the same name
// as those specified in the headers argument. Other headers remain unchanged.
func (r Response) SetHeaders(headers map[string][]string) error {
	arg, err := bridge.WrapHeaders(headers)
	if err != nil {
		return err
	}

	err = r.Ask(`kong.response.set_headers`, arg, nil)
	return err
}

// kong.Response.Exit() interrupts the current processing and produces a response.
// It is typical to see plugins using it to produce a response before Kong
// has a chance to proxy the request (e.g. an authentication plugin rejecting
// a request, or a caching plugin serving a cached response).
//
// This function closes the channel back to Kong main process, so any use of a
// PDK function after this would trigger a run-time panic.  It is recommended to
// stop all processing and return immediately from the handler.
//
// Calling `kong.Response.Exit()` will interrupt the execution flow of
// plugins in the current phase. Subsequent phases will still be invoked.
// E.g. if a plugin called `kong.Response.Exit()` in the `access` phase, no
// other plugin will be executed in that phase, but the `header_filter`,
// `body_filter`, and `log` phases will still be executed, along with their
// plugins. Plugins should thus be programmed defensively against cases when
// a request was **not** proxied to the Service, but instead was produced by
// Kong itself.
//
// The first argument `status` will set the status code of the response that
// will be seen by the client.
//
// The second, `body` argument will set the response body. No special processing
// will be done, and the body will be sent as-is.  It is the caller's responsibility
// to set the appropriate Content-Type header via the third argument.   On gRPC
// we cannot send the `body` with this function at the moment at least, so what it does
// instead is that it sends "body" in `grpc-message` header instead.
//
// The third, `headers` argument can be a table specifying response headers to send.
// If non nil, its behavior is similar to `kong.response.set_headers()`.
//
// Unless manually specified, this method will automatically set the
// Content-Length header in the produced response for convenience.

func (r Response) Exit(status int, body string, headers map[string][]string) {
	h, _ := bridge.WrapHeaders(headers)
	arg := kong_plugin_protocol.ExitArgs{
		Status:  int32(status),
		Body:    body,
		Headers: h,
	}
	r.Ask(`kong.response.exit`, &arg, nil)
	r.Close()
}

// kong.Response.ExitStatus() terminates current processing just like kong.Response.Exit()
// without setting the body or headers.
func (r Response) ExitStatus(status int) {
	arg := kong_plugin_protocol.ExitArgs{
		Status: int32(status),
	}
	r.Ask(`kong.response.exit`, &arg, nil)
	r.Close()
}
