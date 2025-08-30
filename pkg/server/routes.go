package server

import (
	"log"
	"moon/pkg/moon"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MoonTable struct {
	Table []*moon.MoonTableElement
}

type MoonPhaseResponse struct {
	EndDays            float64
	EndDayIllumination float64
	EndDayPhase        string
	EndDayPhaseEmoji   string

	CurrentDays         float64
	CurrentIllumination float64
	CurrentPhase        string
	CurrentPhaseEmoji   string

	BeginDays            float64
	BeginDayIllumination float64
	BeginDayPhase        string
	BeginDayPhaseEmoji   string

	Zodiac string

	FullDaysBegin   float64
	FullDaysCurrent float64
	FullDaysEnd     float64

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

	var beginDuration, currentDuration, endDuration time.Duration
	beginDuration, currentDuration, endDuration, resp.Zodiac = moon.CurrentMoonDays(tGiven, loc)

	resp.FullDaysBegin = beginDuration.Minutes() / moon.Fminute
	resp.FullDaysCurrent = currentDuration.Minutes() / moon.Fminute
	resp.FullDaysEnd = endDuration.Minutes() / moon.Fminute

	resp.BeginDays = toFixed(resp.FullDaysBegin, 2)
	resp.CurrentDays = toFixed(resp.FullDaysCurrent, 2)
	resp.EndDays = toFixed(resp.FullDaysEnd, 2)

	resp.FullIlluminationCurrent, resp.FullIlluminationBeginDay, resp.FullIlluminationEndDay, resp.CurrentPhase, resp.CurrentPhaseEmoji, resp.BeginDayPhase, resp.BeginDayPhaseEmoji, resp.EndDayPhase, resp.EndDayPhaseEmoji = moon.CurrentMoonPhase(tGiven, loc)
	resp.CurrentIllumination = toFixed(resp.FullIlluminationCurrent*100, 2)
	resp.BeginDayIllumination = toFixed(resp.FullIlluminationBeginDay*100, 2)
	resp.EndDayIllumination = toFixed(resp.FullIlluminationEndDay*100, 2)

	return c.JSON(resp)
}

func (s *Server) getCurrentMoonTableV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	resp := MoonTable{}
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)

	resp.Table = moon.GenerateMoonTable(tGiven)

	return c.JSON(resp.Table)
}
