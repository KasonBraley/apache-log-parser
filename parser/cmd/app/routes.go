package main

func (s *server) routes() {
	s.router.Handle("/upload", s.handleUploadPost())
}
