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

type Request struct {
	Method  string
	Url     string
	Headers map[string][]string
	Body    string
}

type Response struct {
	Status  int
	Message string
	Headers map[string][]string
	Body    string
}

type clientEnv struct {
	t   *testing.T
	req Request
	res Response
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

func (e clientEnv) Errorf(format string, args ...interface{}) {
	e.t.Errorf(format, args...)
}

func (e clientEnv) Handle(method string, args_d []byte) []byte {
	switch method {
	case "kong.response.set_header":
		{
			out := new(kong_plugin_protocol.KV)
			err := proto.Unmarshal(args_d, out)
			if err != nil {
				e.t.Error(err)
			}
			e.res.Headers[out.K] = []string{out.V.GetStringValue()}
			return nil
		}

	default:
		{
			e.t.Errorf("unknown method: \"%v\"", method)
			return nil
		}
	}
}

func mockPdk(e clientEnv) *pdk.PDK {
	b := bridge.New(bridgetest.MockFunc(e))
	return &pdk.PDK{
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
}

func (req Request) Exec(t *testing.T, config interface{}) (res Response, err error) {
	err = req.Validate()
	if err != nil {
		return
	}

	res = Response{
		Status:  200,
		Headers: make(map[string][]string),
	}

	mPdk := mockPdk(clientEnv{t, req, res})

	if h, ok := config.(interface{ Certificate(*pdk.PDK) }); ok {
		t.Log("Certificate")
		h.Certificate(mPdk)
	}
	if h, ok := config.(interface{ Rewrite(*pdk.PDK) }); ok {
		t.Log("Rewrite")
		h.Rewrite(mPdk)
	}
	if h, ok := config.(interface{ Access(*pdk.PDK) }); ok {
		t.Log("Access")
		h.Access(mPdk)
	}
	if h, ok := config.(interface{ Response(*pdk.PDK) }); ok {
		t.Log("Response")
		h.Response(mPdk)
	}
	if h, ok := config.(interface{ Log(*pdk.PDK) }); ok {
		t.Log("Log")
		h.Log(mPdk)
	}

	return
}
