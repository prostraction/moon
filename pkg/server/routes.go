package server

import (
	"log"
	"moon/pkg/moon"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

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

	resp.EndDay = new(MoonStat)
	resp.CurrentState = new(MoonStat)
	resp.BeginDay = new(MoonStat)
	resp.Info = new(FullInfo)

	resp.Info.MoonDaysBegin = beginDuration.Minutes() / moon.Fminute
	resp.Info.MoonDaysCurrent = currentDuration.Minutes() / moon.Fminute
	resp.Info.MoonDaysEnd = endDuration.Minutes() / moon.Fminute

	resp.BeginDay.MoonDays = toFixed(resp.Info.MoonDaysBegin, 2)
	resp.CurrentState.MoonDays = toFixed(resp.Info.MoonDaysCurrent, 2)
	resp.EndDay.MoonDays = toFixed(resp.Info.MoonDaysEnd, 2)

	resp.Info.IlluminationCurrent, resp.Info.IlluminationBeginDay, resp.Info.IlluminationEndDay, resp.CurrentDay.Phase, resp.CurrentDay.PhaseEmoji, resp.BeginDay.Phase, resp.BeginDay.PhaseEmoji, resp.EndDay.Phase, resp.EndDay.PhaseEmoji = moon.CurrentMoonPhase(tGiven, loc)

	resp.BeginDay.Illumination = toFixed(resp.Info.IlluminationBeginDay*100, 2)
	resp.CurrentState.Illumination = toFixed(resp.Info.IlluminationCurrent*100, 2)
	resp.EndDay.Illumination = toFixed(resp.Info.IlluminationEndDay*100, 2)

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
