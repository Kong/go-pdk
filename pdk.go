package pdk

import (
	"github.com/kong/go-pdk/client"
	"github.com/kong/go-pdk/ip"
	"github.com/kong/go-pdk/log"
	"github.com/kong/go-pdk/nginx"
	"github.com/kong/go-pdk/node"
	"github.com/kong/go-pdk/request"
	"github.com/kong/go-pdk/response"
	"github.com/kong/go-pdk/router"
	"github.com/kong/go-pdk/service"
	service_request "github.com/kong/go-pdk/service/request"
	service_response "github.com/kong/go-pdk/service/response"
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

func Init(ch chan string) *PDK {
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
