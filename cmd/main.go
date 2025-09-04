package main

import "moon/pkg/server"

func main() {
	s := server.Server{}
	app := s.NewRouter()
	app.Listen(":9999")
}
