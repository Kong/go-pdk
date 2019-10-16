package service

import (
	"fmt"
)

type Service struct {
	ch chan string
}

func NewService(ch chan string) *Service {
	return &Service{ch: ch}
}

func (s *Service) SetUpstream(host string) {
	s.ch <- fmt.Sprintf(`kong.response.set_upstream:%s`, host)
	_ = <- s.ch
}

func (s *Service) SetTarget(host string, port int) {
	s.ch <- fmt.Sprintf(`kong.response.set_upstream:["%s", "%d"]`, host, port)
	_ = <- s.ch
}

// TODO set_tls_cert_key
