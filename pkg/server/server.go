package server

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
}

func (s *Server) NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "Fiber",
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Get("/v1/getCurrentMoonTable", s.getCurrentMoonTableV1)
	app.Get("/v1/getCurrentMoonPhase", s.getCurrentMoonPhaseV1)
	return app
}
