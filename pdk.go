package pdk

import (
	"github.com/kong/go-pdk/client"
	"github.com/kong/go-pdk/log"
	"github.com/kong/go-pdk/nginx"
	"github.com/kong/go-pdk/request"
	"github.com/kong/go-pdk/response"
	"github.com/kong/go-pdk/router"
)

type PDK struct {
	Client   *client.Client
	Log      *log.Log
	Nginx    *nginx.Nginx
	Request  *request.Request
	Response *response.Response
	Router   *router.Router
}

func Init(ch chan string) *PDK {
	return &PDK{
		Client:   client.NewClient(ch),
		Log:      log.NewLog(ch),
		Nginx:    nginx.NewNginx(ch),
		Request:  request.NewRequest(ch),
		Response: response.NewResponse(ch),
		Router:   router.NewRouter(ch),
	}
}
