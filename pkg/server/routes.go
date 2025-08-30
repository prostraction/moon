package server

import (
	"log"
	"moon/pkg/moon"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MoonPhaseResponse struct {
	Days             float64
	DaysFull         float64
	Illumination     float64
	IlluminationFull float64
	Phase            string
	PhaseEmoji       string
	Zodiac           string
}

func (s *Server) getCurrentMoonPhaseV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	resp := MoonPhaseResponse{}
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)

	var duration time.Duration
	duration, resp.Zodiac = moon.CurrentMoonDays(tGiven)
	resp.DaysFull = duration.Minutes()
	resp.DaysFull = resp.DaysFull / 60 / 24
	resp.Days = toFixed(resp.DaysFull, 2)

	resp.IlluminationFull, resp.Phase, resp.PhaseEmoji = moon.CurrentMoonPhase(tGiven)
	resp.Illumination = toFixed(resp.IlluminationFull*100, 2)

	return c.JSON(resp)
}
