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

	BeginDayIllumination float64
	BeginDayPhase        string
	BeginDayPhaseEmoji   string

	EndDayIllumination float64
	EndDayPhase        string
	EndDayPhaseEmoji   string

	Zodiac                   string
	FullDays                 float64
	FullIlluminationCurrent  float64
	FullIlluminationBeginDay float64
	FullIlluminationEndDay   float64
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

	resp.FullIlluminationCurrent, resp.FullIlluminationBeginDay, resp.FullIlluminationEndDay, resp.CurrentPhase, resp.CurrentPhaseEmoji, resp.BeginDayPhase, resp.BeginDayPhaseEmoji, resp.EndDayPhase, resp.EndDayPhaseEmoji = moon.CurrentMoonPhase(tGiven, loc)
	resp.CurrentIllumination = toFixed(resp.FullIlluminationCurrent*100, 2)
	resp.BeginDayIllumination = toFixed(resp.FullIlluminationBeginDay*100, 2)
	resp.EndDayIllumination = toFixed(resp.FullIlluminationEndDay*100, 2)

	return c.JSON(resp)
}
