package server

import (
	"moon/pkg/moon"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/v1/moonTableCurrent", s.moonTableCurrentV1)
	app.Get("/v1/moonTableYear", s.moonTableYearV1)

	app.Get("/v1/moonPhaseCurrent", s.moonPhaseCurrentV1)
	app.Get("/v1/moonPhaseTimestamp", s.moonPhaseTimestampV1)
	app.Get("/v1/moonPhaseDate", s.moonPhaseDatetV1)

	s.moonCache = new(moon.Cache)

	return app
}
