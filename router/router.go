package router

import (
	"encoding/json"
	"github.com/kong/go-pdk/entities"
)

type Router struct {
	ch chan string
}

func NewRouter(ch chan string) *Router {
	return &Router{ch: ch}
}

func (c *Router) GetRoute() *entities.Route {
	c.ch <- `kong.router.get_route`
	reply := <-c.ch
	if reply == "null" {
		return nil
	}
	route := entities.Route{}
	json.Unmarshal([]byte(reply), &route)
	return &route
}

func (c *Router) GetService() *entities.Service {
	c.ch <- `kong.router.get_service`
	reply := <-c.ch
	if reply == "null" {
		return nil
	}
	service := entities.Service{}
	json.Unmarshal([]byte(reply), &service)
	return &service
}
