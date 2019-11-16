package router

import (
	"testing"

	"github.com/Kong/go-pdk/entities"
	"github.com/stretchr/testify/assert"
)

var router Router
var ch chan string

func init() {
	ch = make(chan string)
	router = New(ch)
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func TestGetRoute(t *testing.T) {
	assert.Equal(t, getName(func() { router.GetRoute() }), "kong.router.get_route:null")

	res := make(chan *entities.Route)
	go func(res chan *entities.Route) { r, _ := router.GetRoute(); res <- r }(res)
	_ = <-ch
	ch <- `
		{
			"id": "foo_id",
			"created_at": 123456,
			"paths": ["/", "/foo", "/bar"]
		}`
	route := <-res
	assert.Equal(t, route.Id, "foo_id")
	assert.Equal(t, route.CreatedAt, 123456)
	assert.Equal(t, *route.Paths, []string{"/", "/foo", "/bar"})
}

func TestGetService(t *testing.T) {
	assert.Equal(t, getName(func() { router.GetService() }), "kong.router.get_service:null")

	res := make(chan *entities.Service)
	go func(res chan *entities.Service) { r, _ := router.GetService(); res <- r }(res)
	_ = <-ch
	ch <- `
		{
			"id": "foo_id",
			"created_at": 123456,
			"path": "/foo",
			"port": 80,
			"host": "example.test",
			"protocol": "http"
		}`
	service := <-res
	assert.Equal(t, service.Id, "foo_id")
	assert.Equal(t, service.CreatedAt, 123456)
	assert.Equal(t, service.Path, "/foo")
	assert.Equal(t, service.Port, 80)
	assert.Equal(t, service.Host, "example.test")
	assert.Equal(t, service.Protocol, "http")
}
