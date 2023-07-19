package robin

import (
	"errors"
)

type SimpleRoundRobin struct {
	Servers      []string
	CurrentIndex int
}

func (s *SimpleRoundRobin) Add(servers ...string) error {
	if len(servers) == 0 {
		return errors.New("must have 1 server at least")
	}

	s.Servers = append(s.Servers, servers...)
	return nil
}

func (s *SimpleRoundRobin) GetPeer() string {
	if len(s.Servers) == 0 {
		return ""
	}

	peer := s.Servers[s.CurrentIndex]

	s.CurrentIndex = (s.CurrentIndex + 1) % len(s.Servers)

	return peer
}
