package server

import (
	"log"
	"moon/pkg/moon"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MoonPhaseResponse struct {
	Days                  float64
	Illumination          float64
	IlluminationDaily     float64
	Phase                 string
	PhaseEmoji            string
	Zodiac                string
	FullDays              float64
	FullIllumination      float64
	FullIlluminationDaily float64
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
	resp.FullDays = duration.Minutes()
	resp.FullDays = resp.FullDays / 60 / 24
	resp.Days = toFixed(resp.FullDays, 2)

	resp.FullIllumination, resp.FullIlluminationDaily, resp.Phase, resp.PhaseEmoji = moon.CurrentMoonPhase(tGiven)
	resp.Illumination = toFixed(resp.FullIllumination*100, 2)
	resp.IlluminationDaily = toFixed(resp.FullIlluminationDaily*100, 2)

	return c.JSON(resp)
}
