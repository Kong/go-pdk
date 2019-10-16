package pdk

import (
	"github.com/kong/go-pdk/client"
	"github.com/kong/go-pdk/log"
	"github.com/kong/go-pdk/nginx"
	"github.com/kong/go-pdk/request"
	"github.com/kong/go-pdk/response"
	"github.com/kong/go-pdk/router"
	"github.com/kong/go-pdk/ip"
	"github.com/kong/go-pdk/node"
	"github.com/kong/go-pdk/service"
	service_request "github.com/kong/go-pdk/service/request"
	service_response "github.com/kong/go-pdk/service/response"
)

type PDK struct {
	Client   *client.Client
	Log      *log.Log
	Nginx    *nginx.Nginx
	Request  *request.Request
	Response *response.Response
	Router   *router.Router
	Ip       *ip.Ip
	Node     *node.Node
	Service  *service.Service
	ServiceRequest *service_request.Request
	ServiceResponse *service_response.Response
}

func Init(ch chan string) *PDK {
	return &PDK{
		Client:         client.NewClient(ch),
		Log:            log.NewLog(ch),
		Nginx:          nginx.NewNginx(ch),
		Request:        request.NewRequest(ch),
		Response:       response.NewResponse(ch),
		Router:         router.NewRouter(ch),
		Ip:             ip.NewIp(ch),
		Node:           node.NewNode(ch),
		Service:        service.NewService(ch),
		ServiceRequest: service_request.NewRequest(ch),
		ServiceResponse: service_response.NewResponse(ch),
	}
}
