package server

import (
	"moon/pkg/moon"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Days         int
	Illumination string
}

func (s *Server) NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "Fiber",
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Get("/getMoonPhase", s.getMoonPhase)
	return app
}

func (s *Server) getMoonPhase(c *fiber.Ctx) error {
	D := c.Query("d", "default")
	DInt, err := strconv.Atoi(D)
	if err != nil {
		return err
	}

	M := c.Query("m", "default")
	MInt, err := strconv.Atoi(M)
	if err != nil {
		return err
	}

	Y := c.Query("y", "default")
	yInt, err := strconv.Atoi(Y)
	if err != nil {
		return err
	}

	L := moon.CalcMoonNumber(yInt)
	LCalc := ((L * 11) - 14) % 30

	s.Days = (LCalc + DInt + MInt) % 30
	//s.Illumination = strconv.Itoa((((LCalc + DInt + MInt) % 30) / 4) % 4)
	return c.JSON(s)
}
