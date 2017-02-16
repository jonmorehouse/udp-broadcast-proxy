package main

import "sync"

type server struct {
	wg     sync.WaitGroup
	stopCh chan (struct{})
}

func newServer(opts *options) (*server, error) {

	return nil, nil
}

func (s *server) start() error {
	return nil
}

func (s *server) stop() error {
	return nil
}
