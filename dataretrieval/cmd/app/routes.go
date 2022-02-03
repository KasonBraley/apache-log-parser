package main

func (s *server) routes() {
	s.router.HandleFunc("/retrieve", s.handleDataGet())
}
