package service

import (
	"fmt"
	"github.com/kong/go-pdk/bridge"
)

type Service struct {
	bridge.PdkBridge
}

func New(ch chan string) *Service {
	return &Service{*bridge.New(ch)}
}

func (s *Service) SetUpstream(host string) {
	_ = s.Ask(fmt.Sprintf(`kong.service.set_upstream:%s`, host))
}

func (s *Service) SetTarget(host string, port int) {
	_ = s.Ask(fmt.Sprintf(`kong.service.set_target:["%s", %d]`, host, port))
}

// TODO set_tls_cert_key
