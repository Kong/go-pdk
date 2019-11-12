package pdk

import (
	"github.com/Kong/go-pdk/client"
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
)

type PDK struct {
	Client          client.Client
	Log             log.Log
	Nginx           nginx.Nginx
	Request         request.Request
	Response        response.Response
	Router          router.Router
	Ip              ip.Ip
	Node            node.Node
	Service         service.Service
	ServiceRequest  service_request.Request
	ServiceResponse service_response.Response
}

func Init(ch chan interface{}) *PDK {
	return &PDK{
		Client:          client.New(ch),
		Log:             log.New(ch),
		Nginx:           nginx.New(ch),
		Request:         request.New(ch),
		Response:        response.New(ch),
		Router:          router.New(ch),
		Ip:              ip.New(ch),
		Node:            node.New(ch),
		Service:         service.New(ch),
		ServiceRequest:  service_request.New(ch),
		ServiceResponse: service_response.New(ch),
	}
}
