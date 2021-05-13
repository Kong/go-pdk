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
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Request struct {
	Method  string
	Url     string
	Headers http.Header
	Body    string
}

func (req Request) clone() Request {
	return Request{
		Method:  req.Method,
		Url:     req.Url,
		Headers: req.Headers.Clone(),
		Body:    req.Body,
	}
}

func (req Request) Validate() error {
	_, err := url.Parse(req.Url)
	if err != nil {
		return err
	}

	if req.Method == "GET" {
		if req.Body != "" {
			return fmt.Errorf("GET requests must not have body, found \"%v\"", req.Body)
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

func (req Request) getForwardedUrl() (*url.URL, error) {
	u := req.Headers.Get("X-Forwarded-Proto")
	if u == "" {
		u = req.Url
	}
	return url.Parse(u)
}

func (req Request) ToResponse() Response {
	return Response{
		Status:  200,
		Message: "OK",
		Headers: req.Headers.Clone(),
		Body:    req.Body,
	}
}

type Response struct {
	Status  int
	Message string
	Headers http.Header
	Body    string
}

func (res Response) clone() Response {
	return Response{
		Status:  res.Status,
		Message: res.Message,
		Headers: res.Headers.Clone(),
		Body:    res.Body,
	}
}

type testEnv struct {
	t          *testing.T
	pdk        *pdk.PDK
	ClientReq  Request
	ServiceReq Request
	ServiceRes Response
	ClientRes  Response
}

func New(t *testing.T, req Request) (env *testEnv, err error) {
	err = req.Validate()
	if err != nil {
		return
	}

	env = &testEnv{
		t:          t,
		ClientReq:  req,
		ServiceReq: req.clone(),
		ServiceRes: Response{},
		ClientRes:  Response{},
	}

	b := bridge.New(bridgetest.MockFunc(env)) // check
	env.pdk = &pdk.PDK{
		Client:          client.Client{b},
		Ctx:             ctx.Ctx{b},
		Log:             log.Log{b},
		Nginx:           nginx.Nginx{b},
		Request:         request.Request{b},
		Response:        response.Response{b},
		Router:          router.Router{b},
		IP:              ip.Ip{b},
		Node:            node.Node{b},
		Service:         service.Service{b},
		ServiceRequest:  service_request.Request{b},
		ServiceResponse: service_response.Response{b},
	}
	return
}

func mergeHeaders(h http.Header, in map[string][]string) {
	for k, l := range in {
		h.Del(k)
		for _, v := range l {
			h.Add(k, v)
		}
	}
}

func (e testEnv) noErr(err error) {
	if err != nil {
		e.t.Error(err)
	}
}

func (e testEnv) Errorf(format string, args ...interface{}) {
	e.t.Errorf(format, args...)
}

func (e *testEnv) Handle(method string, args_d []byte) []byte {
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
		u, err := e.ClientReq.getForwardedUrl()
		e.noErr(err)
		out = bridge.WrapString(u.Scheme)

	case "kong.request.get_forwarded_host":
		u, err := e.ClientReq.getForwardedUrl()
		e.noErr(err)
		out = bridge.WrapString(u.Hostname())

	case "kong.request.get_forwarded_port":
		u, err := e.ClientReq.getForwardedUrl()
		e.noErr(err)
		out = &kong_plugin_protocol.Int{V: getPort(u)}

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

	case "kong.request.get_header":
		args := kong_plugin_protocol.String{}
		e.noErr(proto.Unmarshal(args_d, &args))
		out = bridge.WrapString(e.ClientReq.Headers.Get(args.V))

	case "kong.request.get_raw_body":
		out = bridge.WrapString(e.ClientReq.Body)

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
			e.ClientRes.Body = args.Body
			headers := bridge.UnwrapHeaders(args.Headers)
			mergeHeaders(e.ClientRes.Headers, headers)
		}

	case "kong.router.get_route":
		out = &kong_plugin_protocol.Route{
			Id: "001:002",
			Name: "route_66",
			Protocols: []string{"http", "tcp"},
			Paths: []string{"/v0/left", "/v1/this"},
		}

	case "kong.router.get_service":
		out = &kong_plugin_protocol.Service{
			Id: "003:004",
			Name: "self_service",
			Protocol: "http",
			Path: "/v0/left",
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
		args := kong_plugin_protocol.String{}
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
		out = bridge.WrapString(e.ServiceRes.Body)

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

func (e *testEnv) DoCertificate(config interface{}) {
	if h, ok := config.(interface{ Certificate(*pdk.PDK) }); ok {
		e.t.Log("Certificate")
		h.Certificate(e.pdk)
	}
}

func (e *testEnv) DoRewrite(config interface{}) {
	if h, ok := config.(interface{ Rewrite(*pdk.PDK) }); ok {
		e.t.Log("Rewrite")
		h.Rewrite(e.pdk)
	}
}

func (e *testEnv) DoAccess(config interface{}) {
	if h, ok := config.(interface{ Access(*pdk.PDK) }); ok {
		e.t.Log("Access")
		h.Access(e.pdk)
	}
}

func (e *testEnv) DoResponse(config interface{}) {
	e.ClientRes = e.ServiceRes.clone()
	if h, ok := config.(interface{ Response(*pdk.PDK) }); ok {
		e.t.Log("Response")
		h.Response(e.pdk)
	}
}

func (e *testEnv) DoPreread(config interface{}) {
	if h, ok := config.(interface{ Preread(*pdk.PDK) }); ok {
		e.t.Log("Preread")
		h.Preread(e.pdk)
	}
}

func (e *testEnv) DoLog(config interface{}) {
	if h, ok := config.(interface{ Log(*pdk.PDK) }); ok {
		e.t.Log("Log")
		h.Log(e.pdk)
	}
}

func (e *testEnv) DoHttp(config interface{}) {
	e.DoAccess(config)
	e.ServiceRes = e.ServiceReq.ToResponse() // assuming an "echo service"
	e.DoResponse(config)
	e.DoLog(config)
}

func (e *testEnv) DoHttps(config interface{}) {
	e.DoCertificate(config)
	e.DoHttp(config)
}

func (e *testEnv) DoStream(config interface{}) {
	e.DoPreread(config)
	e.DoLog(config)
}

func (e *testEnv) DoTLS(config interface{}) {
	e.DoCertificate(config)
	e.DoStream(config)
}
