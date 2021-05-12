package test

import (
	"fmt"
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
	// 	"google.golang.org/protobuf/types/known/structpb"
)

type headersMap map[string][]string

func (h headersMap) cloneTo(out_h map[string][]string) {
	for k, v := range h {
		out_v := make([]string, len(v))
		for i, sub_v := range v {
			out_v[i] = sub_v
		}
		out_h[k] = out_v
	}
}
func (h headersMap) clone() map[string][]string {
	out_h := make(map[string][]string, len(h))
	h.cloneTo(out_h)
	return out_h
}

type Request struct {
	Method  string
	Url     string
	Headers map[string][]string
	Body    string
}

func (req Request) Validate() error {
	if req.Method == "GET" {
		if req.Body != "" {
			return fmt.Errorf("GET requests must not have body, found \"%v\"", req.Body)
		}
		return nil
	}
	return fmt.Errorf("Unsupported method \"%v\"", req.Method)
}

func (req Request) clone() Request {
	return Request{
		Method:  req.Method,
		Url:     req.Url,
		Headers: headersMap(req.Headers).clone(),
		Body:    req.Body,
	}
}

func (req Request) ToResponse(res *Response) {
	res.Status = 200
	res.Message = "OK"
	headersMap(req.Headers).cloneTo(res.Headers)
	res.Body = req.Body
}

type Response struct {
	Status  int
	Message string
	Headers map[string][]string
	Body    string
}

func (res Response) cloneTo(out *Response) {
	out.Status = res.Status
	out.Message = res.Message
	headersMap(res.Headers).cloneTo(out.Headers)
	out.Body = res.Body
}

func (res Response) clone() Response {
	out := Response{Headers: map[string][]string{}}
	res.cloneTo(&out)
	return out
}

type testEnv struct {
	t          *testing.T
	pdk        *pdk.PDK
	ClientReq  Request
	ServiceReq Request
	ServiceRes Response
	ClientRes  Response
}

func New(t *testing.T, req Request) (env testEnv, err error) {
	if req.Headers == nil {
		req.Headers = map[string][]string{}
	}
	err = req.Validate()
	if err != nil {
		return
	}

	env = testEnv{
		t:          t,
		ClientReq:  req,
		ServiceReq: req.clone(),
		ServiceRes: Response{
			Headers: map[string][]string{},
		},
		ClientRes: Response{
			Headers: map[string][]string{},
		},
	}

	b := bridge.New(bridgetest.MockFunc(&env)) // check
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

	case "kong.request.get_headers":
		out, err = bridge.WrapHeaders(e.ClientReq.Headers)

	case "kong.client.get_credential":
		out = &kong_plugin_protocol.AuthenticatedCredential{Id: "000:00", ConsumerId: "000:01"}

	case "kong.client.load_consumer", "kong.client.get_consumer":
		out = &kong_plugin_protocol.Consumer{Id: "001", Username: "Jon Doe"}

	case "kong.client.authenticate":
		// accept anything

	case "kong.client.get_protocol":
		out = bridge.WrapString("https")

	case "kong.response.set_header":
		{
			args := new(kong_plugin_protocol.KV)
			e.noErr(proto.Unmarshal(args_d, args))
			e.ClientRes.Headers[args.K] = []string{args.V.GetStringValue()}
			return nil
		}

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
	e.ServiceRes.cloneTo(&e.ClientRes)
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
	e.ServiceReq.ToResponse(&e.ServiceRes) // assuming an "echo service"
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
