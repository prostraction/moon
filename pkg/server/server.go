package server

import (
	"moon/pkg/moon"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	moonCache *moon.Cache
}

func (s *Server) NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "Fiber",
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Get("/v1/moonTableCurrent", s.moonTableCurrentV1)
	app.Get("/v1/moonTableYear", s.moonTableYearV1)

	app.Get("/v1/moonPhaseCurrent", s.moonPhaseCurrentV1)
	app.Get("/v1/moonPhaseTimestamp", s.moonPhaseTimestampV1)
	app.Get("/v1/moonPhaseDate", s.moonPhaseDatetV1)

	s.moonCache = new(moon.Cache)

	return app
}
