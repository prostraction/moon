package server

import (
	"log"
	"moon/pkg/moon"
	"strconv"
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

/*    MOON PHASE    */

func (s *Server) moonPhaseCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	return s.moonPhaseV1(c, tGiven)
}

func (s *Server) moonPhaseTimestampV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	tStr := c.Query("t", "0")
	t, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		t = time.Now().Unix()
	}
	tm := time.Unix(t, 0)
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, loc)
	return s.moonPhaseV1(c, tGiven)
}

func (s *Server) moonPhaseDatetV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	year := strToInt(c.Query("year", "1970"), 1970, 0)
	month := strToInt(c.Query("month", "1"), 1, 12)
	day := strToInt(c.Query("day", "1"), 1, 31)

	hour := strToInt(c.Query("hour", "0"), 0, 23)
	minute := strToInt(c.Query("minute", "0"), 0, 59)
	second := strToInt(c.Query("second", "0"), 0, 59)

	tGiven := time.Date(year, moon.GetMonth(month), day, hour, minute, second, 0, loc)
	return s.moonPhaseV1(c, tGiven)
}

func (s *Server) moonPhaseV1(c *fiber.Ctx, tGiven time.Time) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	resp := MoonPhaseResponse{}
	var beginDuration, currentDuration, endDuration time.Duration
	beginDuration, currentDuration, endDuration, resp.Zodiac = s.moonCache.CurrentMoonDays(tGiven, loc)

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

/*    MOON TABLE    */

func (s *Server) moonTableCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableYearV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, err := moon.SetTimezoneLocFromString(utc)
	if err != nil {
		log.Println(err)
	}

	yearStr := c.Query("year", "1970")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = time.Now().Year()
	}

	tGiven := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableV1(c *fiber.Ctx, tGiven time.Time) error {
	resp := MoonTable{}
	resp.Table = s.moonCache.GenerateMoonTable(tGiven)
	return c.JSON(resp.Table)
}
