/*
Utilities to test plugins.

Trivial example:

	package main

	import (
		"testing"

		"github.com/Kong/go-pdk/test"
		"github.com/stretchr/testify/assert"
	)

	func TestPlugin(t *testing.T) {
		env, err := test.New(t, test.Request{
			Method: "GET",
			Url:    "http://example.com?q=search&x=9",
			Headers: map[string][]string{ "X-Hi": {"hello"}, },
		})
		assert.NoError(t, err)

		env.DoHttps(&Config{})
		assert.Equal(t, 200, env.ClientRes.Status)
		assert.Equal(t, "Go says Hi!", env.ClientRes.Headers.Get("x-hello-from-go"))
	}

in short:

1. Create a test environment passing a test.Request{} object to the test.New() function.

2. Create a Config{} object (or the appropriate config structure of the plugin)

3. env.DoHttps(t, &config) will pass the request object through the plugin, exercising
each event handler and return (if there's no hard error) a simulated response object.
There are other env.DoXXX(t, &config) functions for HTTP, TCP, TLS and individual phases.

3.5 The http and https functions assume the service response will be an "echo" of the
request (same body and headers) if you need a different service response, use the
individual phase methods and set the env.ServiceRes object manually.

4. Do assertions to verify the service request and client response are as expected.
*/
package test

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/client"
	"github.com/Kong/go-pdk/ctx"
	"github.com/Kong/go-pdk/ip"
	"github.com/Kong/go-pdk/log"
	"github.com/Kong/go-pdk/nginx"
	"github.com/Kong/go-pdk/node"
	"github.com/Kong/go-pdk/request"
	"github.com/Kong/go-pdk/response"
	"github.com/Kong/go-pdk/router"
	"github.com/Kong/go-pdk/service"
	service_request "github.com/Kong/go-pdk/service/request"
	service_response "github.com/Kong/go-pdk/service/response"

	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// The Request type represents the request received from the client
// or the request sent to the service.
type Request struct {
	Method  string
	Url     string
	Headers http.Header
	Body    []byte
}

func (req Request) clone() Request {
	return Request{
		Method:  req.Method,
		Url:     req.Url,
		Headers: req.Headers.Clone(),
		Body:    req.Body,
	}
}

func mergeHeaders(h http.Header, in http.Header) http.Header {
	for k, l := range in {
		h.Del(k)
		for _, v := range l {
			h.Add(k, v)
		}
	}
	return h
}

// Validate verifies a request and normalizes the headers.
// (to make them case-insensitive)
func (req *Request) Validate() error {
	_, err := url.Parse(req.Url)
	if err != nil {
		return err
	}

	req.Headers = mergeHeaders(make(http.Header), req.Headers)

	if req.Method == "GET" {
		if len(req.Body) != 0 {
			return fmt.Errorf("GET requests must not have body, found \"%v\"", req.Body)
		}
		return nil
	}
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		if len(req.Body) == 0 {
			return fmt.Errorf("%s requests must have body, found \"%v\"", req.Method, req.Body)
		}
		return nil
	}
	return fmt.Errorf("Unsupported method \"%v\"", req.Method)
}

func getPort(u *url.URL) int32 {
	p := u.Port()
	if p == "" {
		if u.Scheme == "https" {
			return 443
		}
		return 80
	}
	portnum, _ := strconv.Atoi(p)
	return int32(portnum)
}

// ToResponse creates a new Response object from a Request,
// simulating an "echo" service.
func (req Request) ToResponse() Response {
	return Response{
		Status:  200,
		Message: "OK",
		Headers: req.Headers.Clone(),
		Body:    req.Body,
	}
}

// The Response type represents the response returned from
// the service or sent to the client.
type Response struct {
	Status  int
	Message string
	Headers http.Header
	Body    []byte
}

func (res *Response) merge(other Response) {
	if other.Status != 0 {
		res.Status = other.Status
	}
	if other.Message != "" {
		res.Message = other.Message
	}
	if other.Headers != nil {
		if res.Headers != nil {
			mergeHeaders(res.Headers, other.Headers)
		} else {
			res.Headers = other.Headers.Clone()
		}
	}
	if len(other.Body) != 0 {
		res.Body = other.Body
	}
}

type envState int

const (
	running envState = iota
	finished
)

type TestEnv struct {
	t           *testing.T
	state       envState
	stateChange chan<- string
	pdk         *pdk.PDK
	ClientReq   Request
	ServiceReq  Request
	ServiceRes  Response
	ClientRes   Response
}

// New creates a new test environment.
func New(t *testing.T, req Request) (env *TestEnv, err error) {
	err = req.Validate()
	if err != nil {
		return
	}

	env = &TestEnv{
		t:          t,
		state:      running,
		ClientReq:  req,
		ServiceReq: req.clone(),
		ServiceRes: Response{Headers: make(http.Header)},
		ClientRes:  Response{Headers: make(http.Header)},
	}

	b := bridge.New(bridgetest.MockFunc(env)) // check
	env.pdk = &pdk.PDK{
		Client:          client.Client{PdkBridge: b},
		Ctx:             ctx.Ctx{PdkBridge: b},
		Log:             log.Log{PdkBridge: b},
		Nginx:           nginx.Nginx{PdkBridge: b},
		Request:         request.Request{PdkBridge: b},
		Response:        response.Response{PdkBridge: b},
		Router:          router.Router{PdkBridge: b},
		IP:              ip.Ip{PdkBridge: b},
		Node:            node.Node{PdkBridge: b},
		Service:         service.Service{PdkBridge: b},
		ServiceRequest:  service_request.Request{PdkBridge: b},
		ServiceResponse: service_response.Response{PdkBridge: b},
	}
	return
}

func (e *TestEnv) noErr(err error) {
	if err != nil {
		e.t.Error(err)
	}
}

// Internal use.  Calls the Errof function with the test context.
func (e *TestEnv) Errorf(format string, args ...interface{}) {
	e.t.Errorf(format, args...)
}

func (e *TestEnv) IsRunning() bool {
	return e.state == running
}

func (e *TestEnv) Finish() {
	e.state = finished
	if e.stateChange != nil {
		e.stateChange <- "finished"
	}
}

func (e *TestEnv) SubscribeStatusChange(ch chan<- string) {
	e.stateChange = ch
}

// Internal use.  Handles a PDK request from the plugin under test.
func (e *TestEnv) Handle(method string, args_d []byte) []byte {
	var out proto.Message
	var err error

	switch method {

	case "kong.client.get_ip", "kong.client.get_forwarded_ip":
		out = bridge.WrapString("10.10.10.1")

	case "kong.client.get_port", "kong.client.get_forwarded_port":
		out = &kong_plugin_protocol.Int{V: 443}

	case "kong.client.get_credential":
		out = &kong_plugin_protocol.AuthenticatedCredential{Id: "000:00", ConsumerId: "000:01"}

	case "kong.client.load_consumer", "kong.client.get_consumer":
		out = &kong_plugin_protocol.Consumer{Id: "001", Username: "Jon Doe"}

	case "kong.client.authenticate":
		// accept anything

	case "kong.client.get_protocol":
		out = bridge.WrapString("https")

	case "kong.ip.is_trusted":
		out = &kong_plugin_protocol.Bool{V: true}

	case "kong.log.alert", "kong.log.crit", "kong.log.err", "kong.log.warn",
		"kong.log.notice", "kong.log.info", "kong.log.debug":
		args := new(structpb.ListValue)
		e.noErr(proto.Unmarshal(args_d, args))
		e.t.Logf("Log (%s): %v", method[strings.LastIndex(method, ".")+1:], args.AsSlice())

	case "kong.node.get_id":
		out = bridge.WrapString("a9777ac2-57e6-482b-a3c4-ef3d6ca41a1f")

	case "kong.node.get_memory_stats":
		out = &kong_plugin_protocol.MemoryStats{
			LuaSharedDicts: &kong_plugin_protocol.MemoryStats_LuaSharedDicts{
				Kong: &kong_plugin_protocol.MemoryStats_LuaSharedDicts_DictStats{
					AllocatedSlabs: 1027,
					Capacity:       4423543,
				},
				KongDbCache: &kong_plugin_protocol.MemoryStats_LuaSharedDicts_DictStats{
					AllocatedSlabs: 4093,
					Capacity:       3424875,
				},
			},
			WorkersLuaVms: []*kong_plugin_protocol.MemoryStats_WorkerLuaVm{
				{HttpAllocatedGc: 123456, Pid: 543},
				{HttpAllocatedGc: 345678, Pid: 876},
			},
		}

	case "kong.request.get_scheme":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = bridge.WrapString(u.Scheme)

	case "kong.request.get_host":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = bridge.WrapString(u.Hostname())

	case "kong.request.get_port":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = &kong_plugin_protocol.Int{V: getPort(u)}

	case "kong.request.get_forwarded_scheme":
		scheme := e.ClientReq.Headers.Get("X-Forwarded-Proto")
		if scheme == "" {
			u, err := url.Parse(e.ClientReq.Url)
			e.noErr(err)
			scheme = u.Scheme
		}
		out = bridge.WrapString(scheme)

	case "kong.request.get_forwarded_host":
		host := e.ClientReq.Headers.Get("X-Forwarded-Host")
		if host == "" {
			u, err := url.Parse(e.ClientReq.Url)
			e.noErr(err)
			host = u.Hostname()
		}
		out = bridge.WrapString(host)

	case "kong.request.get_forwarded_port":
		port := e.ClientReq.Headers.Get("X-Forwarded-Port")
		if port != "" {
			p, err := strconv.Atoi(port)
			e.noErr(err)
			out = &kong_plugin_protocol.Int{V: int32(p)}
		} else {
			u, err := url.Parse(e.ClientReq.Url)
			e.noErr(err)
			out = &kong_plugin_protocol.Int{V: getPort(u)}
		}

	case "kong.request.get_http_version":
		out = &kong_plugin_protocol.Number{V: 1.1}

	case "kong.request.get_method":
		out = bridge.WrapString(e.ClientReq.Method)

	case "kong.request.get_path":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = bridge.WrapString(u.Path)

	case "kong.request.get_path_with_query":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = bridge.WrapString(u.Path + "?" + u.RawQuery)

	case "kong.request.get_raw_query":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = bridge.WrapString(u.RawQuery)

	case "kong.request.get_query_arg":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out = bridge.WrapString(u.Query()[args.V][0])

	case "kong.request.get_query":
		u, err := url.Parse(e.ClientReq.Url)
		e.noErr(err)
		out, err = bridge.WrapHeaders(u.Query())
		e.noErr(err)

	case "kong.request.get_header":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		out = bridge.WrapString(e.ClientReq.Headers.Get(args.V))

	case "kong.request.get_raw_body":
		out = bridge.WrapByteString(e.ClientReq.Body)

	case "kong.request.get_headers":
		out, err = bridge.WrapHeaders(e.ClientReq.Headers)

	case "kong.response.get_status":
		out = &kong_plugin_protocol.Int{V: int32(e.ClientRes.Status)}

	case "kong.response.get_header":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		out = bridge.WrapString(e.ClientRes.Headers.Get(args.V))

	case "kong.response.get_headers":
		out, err = bridge.WrapHeaders(e.ClientRes.Headers)

	case "kong.response.get_source":
		out = bridge.WrapString("service")

	case "kong.response.set_status":
		args := kong_plugin_protocol.Int{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ClientRes.Status = int(args.V)

	case "kong.response.set_header":
		args := kong_plugin_protocol.KV{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ClientRes.Headers.Set(args.K, args.V.GetStringValue())

	case "kong.response.add_header":
		args := kong_plugin_protocol.KV{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ClientRes.Headers.Add(args.K, args.V.GetStringValue())

	case "kong.response.clear_header":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ClientRes.Headers.Del(args.V)

	case "kong.response.set_headers":
		args := structpb.Struct{}
		e.noErr(proto.Unmarshal(args_d, &args))
		headers := bridge.UnwrapHeaders(&args)
		mergeHeaders(e.ClientRes.Headers, headers)

	case "kong.response.exit":
		args := kong_plugin_protocol.ExitArgs{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ClientRes.Status = int(args.Status)
		if args.Headers != nil {
			e.Finish()
			e.ClientRes.Body = args.Body
			headers := bridge.UnwrapHeaders(args.Headers)
			mergeHeaders(e.ClientRes.Headers, headers)
		}
		e.state = finished
		e.ClientRes.Body = args.Body

	case "kong.router.get_route":
		out = &kong_plugin_protocol.Route{
			Id:        "001:002",
			Name:      "route_66",
			Protocols: []string{"http", "tcp"},
			Paths:     []string{"/v0/left", "/v1/this"},
		}

	case "kong.router.get_service":
		out = &kong_plugin_protocol.Service{
			Id:       "003:004",
			Name:     "self_service",
			Protocol: "http",
			Path:     "/v0/left",
		}

	case "kong.service.set_upstream", "kong.service.set_target":
		// no visible effects yet

	case "kong.service.request.set_scheme":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		u, err := url.Parse(e.ServiceReq.Url)
		e.noErr(err)
		u.Scheme = args.V
		e.ServiceReq.Url = u.String()

	case "kong.service.request.set_path":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		u, err := url.Parse(e.ServiceReq.Url)
		e.noErr(err)
		u.Path = args.V
		e.ServiceReq.Url = u.String()

	case "kong.service.request.set_raw_query":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		u, err := url.Parse(e.ServiceReq.Url)
		e.noErr(err)
		u.RawQuery = args.V
		e.ServiceReq.Url = u.String()

	case "kong.service.request.set_method":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ServiceReq.Method = args.V

	case "kong.service.request.set_query":
		args := structpb.Struct{}
		e.noErr(proto.Unmarshal(args_d, &args))
		query := bridge.UnwrapHeaders(&args)
		u, err := url.Parse(e.ServiceReq.Url)
		e.noErr(err)
		u.RawQuery = url.Values(query).Encode()
		e.ServiceReq.Url = u.String()

	case "kong.service.request.set_header":
		args := kong_plugin_protocol.KV{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ServiceReq.Headers.Set(args.K, args.V.GetStringValue())

	case "kong.service.request.add_header":
		args := kong_plugin_protocol.KV{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ServiceReq.Headers.Add(args.K, args.V.GetStringValue())

	case "kong.service.request.clear_header":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ServiceReq.Headers.Del(args.V)

	case "kong.service.request.set_headers":
		args := structpb.Struct{}
		e.noErr(proto.Unmarshal(args_d, &args))
		headers := bridge.UnwrapHeaders(&args)
		mergeHeaders(e.ServiceReq.Headers, headers)

	case "kong.service.request.set_raw_body":
		args := kong_plugin_protocol.ByteString{}
		e.noErr(proto.Unmarshal(args_d, &args))
		e.ServiceRes.Body = args.V

	case "kong.service.response.get_status":
		out = &kong_plugin_protocol.Int{V: int32(e.ServiceRes.Status)}

	case "kong.service.response.get_headers":
		out, err = bridge.WrapHeaders(e.ServiceRes.Headers)

	case "kong.service.response.get_header":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		out = bridge.WrapString(e.ServiceRes.Headers.Get(args.V))

	case "kong.service.response.get_raw_body":
		out = bridge.WrapByteString(e.ServiceRes.Body)

	default:
		e.t.Errorf("unknown method: \"%v\"", method)
	}

	if out != nil {
		e.noErr(err)
		out_d, err := proto.Marshal(out)
		e.noErr(err)
		return out_d
	}
	return nil
}

// DoCertificate tests the Certificate method of the plugin
// with the Request in the test environment and the plugin
// configuration passed in the argument.
func (e *TestEnv) DoCertificate(config interface{}) {
	if !e.IsRunning() {
		return
	}
	if h, ok := config.(interface{ Certificate(*pdk.PDK) }); ok {
		e.t.Log("Certificate")
		h.Certificate(e.pdk)
	}
}

// DoRewrite tests the Rewrite method of the plugin
// with the Request in the test environment and the plugin
// configuration passed in the argument.
func (e *TestEnv) DoRewrite(config interface{}) {
	if !e.IsRunning() {
		return
	}
	if h, ok := config.(interface{ Rewrite(*pdk.PDK) }); ok {
		e.t.Log("Rewrite")
		h.Rewrite(e.pdk)
	}
}

// DoAccess tests the Access method of the plugin
// with the Request in the test environment and the plugin
// configuration passed in the argument.
func (e *TestEnv) DoAccess(config interface{}) {
	if !e.IsRunning() {
		return
	}
	if h, ok := config.(interface{ Access(*pdk.PDK) }); ok {
		e.t.Log("Access")
		h.Access(e.pdk)
	}
}

// DoResponse tests the Response method of the plugin
// with the Request in the test environment and the plugin
// configuration passed in the argument.
// Before calling the plugin, the simulated client response is
// updated with the service response, if any.
func (e *TestEnv) DoResponse(config interface{}) {
	if !e.IsRunning() {
		return
	}
	e.ClientRes.merge(e.ServiceRes)
	if h, ok := config.(interface{ Response(*pdk.PDK) }); ok {
		e.t.Log("Response")
		h.Response(e.pdk)
	}
}

// DoPreread tests the Preread method of a streaming plugin.
func (e *TestEnv) DoPreread(config interface{}) {
	if !e.IsRunning() {
		return
	}
	if h, ok := config.(interface{ Preread(*pdk.PDK) }); ok {
		e.t.Log("Preread")
		h.Preread(e.pdk)
	}
}

// DoLog tests the Log method of the plugin
// with the plugin configuration passed in the argument.
func (e *TestEnv) DoLog(config interface{}) {
	if !e.IsRunning() {
		return
	}
	if h, ok := config.(interface{ Log(*pdk.PDK) }); ok {
		e.t.Log("Log")
		h.Log(e.pdk)
	}
}

// DoHttp simulates an HTTP request/response cycle passing
// through the Rewrite, Access, Response and Log methods
// of the plugin.
//
// Between the Access and Response methods, a service response
// is created from the service request (possibly modified by
// the previous methods), simulating an "echo" service.
// If you need a different kind of service, use the individual
// methods (e.DoRewrite(), e.DoAccess(), e.DoResponse() and e.DoLog())
func (e *TestEnv) DoHttp(config interface{}) {
	e.DoRewrite(config)
	e.DoAccess(config)
	e.ServiceRes = e.ServiceReq.ToResponse() // assuming an "echo service"
	e.DoResponse(config)
	e.DoLog(config)
}

// DoHttps simulates an HTTPS request/response cycle passing
// through the Certificate method and then all the same as the
// DoHttp function.
func (e *TestEnv) DoHttps(config interface{}) {
	e.DoCertificate(config)
	e.DoHttp(config)
}

// DoStream simulates a TCP stream (for streaming plugins).
func (e *TestEnv) DoStream(config interface{}) {
	e.DoPreread(config)
	e.DoLog(config)
}

// DoTLS simulates a TLS stream (for streaming plugins).
func (e *TestEnv) DoTLS(config interface{}) {
	e.DoCertificate(config)
	e.DoStream(config)
}
