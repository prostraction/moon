package main

import "moon/internal/server"

func main() {
	s := server.Server{}
	app := s.NewRouter()
	app.Listen(":9998")
}
