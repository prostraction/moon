package server

import (
	"log"
	"moon/pkg/moon"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MoonPhaseResponse struct {
	Days                float64
	CurrentIllumination float64
	CurrentPhase        string
	CurrentPhaseEmoji   string

	DailyIllumination float64
	DailyPhase        string
	DailyPhaseEmoji   string

	Zodiac                  string
	FullDays                float64
	FullIlluminationCurrent float64
	FullIlluminationDaily   float64
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

	resp.FullIlluminationCurrent, resp.FullIlluminationDaily, resp.CurrentPhase, resp.CurrentPhaseEmoji, resp.DailyPhase, resp.DailyPhaseEmoji = moon.CurrentMoonPhase(tGiven)
	resp.CurrentIllumination = toFixed(resp.FullIlluminationCurrent*100, 2)
	resp.DailyIllumination = toFixed(resp.FullIlluminationDaily*100, 2)

	return c.JSON(resp)
}
