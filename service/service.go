package service

import (
	"github.com/Kong/go-pdk/bridge"
)

type Service struct {
	bridge.PdkBridge
}

func New(ch chan interface{}) Service {
	return Service{bridge.New(ch)}
}

func (s Service) SetUpstream(host string) error {
	_, err := s.Ask(`kong.service.set_upstream`, host)
	return err
}

func (s Service) SetTarget(host string, port int) error {
	_, err := s.Ask(`kong.service.set_target`, host, port)
	return err
}

// TODO set_tls_cert_key
