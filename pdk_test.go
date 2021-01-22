package pdk

import (
	"github.com/Kong/go-pdk/response"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type DemoConfig struct {
	Message string
}

func New() interface{} {
	return &DemoConfig{}
}

func (conf *DemoConfig) Access(kong *PDK) {
	kong.Response.AddHeader("x-test-me", conf.Message)
	kong.Response.ExitStatus(http.StatusForbidden)
}


type responseMock struct {
	response.Response
	mock.Mock
}

func (r responseMock) AddHeader(k string, v string) error {
	args := r.Called(k, v)
	return args.Error(0)
}

func (r responseMock) ExitStatus(code int) {
	r.Called(code)
}


func TestPDKCanBeMockedOut(t *testing.T) {
	// create an instance of our test object
	testResp := new(responseMock)
	// setup expectations
	testResp.On("AddHeader", "x-test-me", "test me").Return(nil)
	testResp.On("ExitStatus", http.StatusForbidden).Return(nil)

	p := &PDK{Response: testResp}

	conf := DemoConfig{Message: "test me"}
	conf.Access(p)
	testResp.AssertExpectations(t)
}