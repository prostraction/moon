package server

import (
	"moon/pkg/moon"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) versionV1(c *fiber.Ctx) error {
	return c.JSON("1.0.3")
}

/*    MOON PHASE    */
func (s *Server) moonPhaseCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := moon.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/
	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	precision := strToInt(c.Query("precision", "2"), 2, 10)
	return s.moonPhaseV1(c, tGiven, precision)
}

func (s *Server) moonPhaseTimestampV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := moon.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tStr := c.Query("t", strconv.FormatInt(time.Now().Unix(), 10))
	t, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		t = time.Now().Unix()
	}
	tm := time.Unix(t, 0)
	tGiven := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, loc)

	precision := strToInt(c.Query("precision", "2"), 2, 10)
	return s.moonPhaseV1(c, tGiven, precision)
}

func (s *Server) moonPhaseDatetV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := moon.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tNow := time.Now()

	year := strToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 0)
	month := strToInt(c.Query("month", strconv.Itoa(int(tNow.Month()))), int(tNow.Month()), 12)
	day := strToInt(c.Query("day", strconv.Itoa(int(tNow.Day()))), int(tNow.Day()), 31)

	hour := strToInt(c.Query("hour", strconv.Itoa(int(tNow.Hour()))), int(tNow.Hour()), 23)
	minute := strToInt(c.Query("minute", strconv.Itoa(int(tNow.Minute()))), int(tNow.Minute()), 59)
	second := strToInt(c.Query("second", strconv.Itoa(int(tNow.Second()))), int(tNow.Second()), 59)

	precision := strToInt(c.Query("precision", "2"), 2, 10)

	tGiven := time.Date(year, moon.GetMonth(month), day, hour, minute, second, 0, loc)
	return s.moonPhaseV1(c, tGiven, precision)
}

func (s *Server) moonPhaseV1(c *fiber.Ctx, tGiven time.Time, precision int) error {
	lang := c.Query("lang", "en")
	utc := c.Query("utc", "UTC:+0")
	loc, _ := moon.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	resp := MoonPhaseResponse{}
	resp.EndDay = new(MoonStat)
	resp.CurrentState = new(MoonStat)
	resp.BeginDay = new(MoonStat)
	resp.info = new(FullInfo)

	var beginDuration, currentDuration, endDuration time.Duration
	beginDuration, currentDuration, endDuration = s.moonCache.CurrentMoonDays(tGiven, loc)

	resp.info.MoonDaysBegin = beginDuration.Minutes() / moon.Fminute
	resp.info.MoonDaysCurrent = currentDuration.Minutes() / moon.Fminute
	resp.info.MoonDaysEnd = endDuration.Minutes() / moon.Fminute

	resp.BeginDay.MoonDays = toFixed(resp.info.MoonDaysBegin, precision)
	resp.CurrentState.MoonDays = toFixed(resp.info.MoonDaysCurrent, precision)
	resp.EndDay.MoonDays = toFixed(resp.info.MoonDaysEnd, precision)

	resp.info.IlluminationCurrent, resp.info.IlluminationBeginDay, resp.info.IlluminationEndDay, resp.CurrentState.Phase, resp.BeginDay.Phase, resp.EndDay.Phase = moon.CurrentMoonPhase(tGiven, loc, lang)

	resp.BeginDay.Illumination = toFixed(resp.info.IlluminationBeginDay*100, precision)
	resp.CurrentState.Illumination = toFixed(resp.info.IlluminationCurrent*100, precision)
	resp.EndDay.Illumination = toFixed(resp.info.IlluminationEndDay*100, precision)

	resp.ZodiacDetailed, resp.BeginDay.Zodiac, resp.CurrentState.Zodiac, resp.EndDay.Zodiac = s.moonCache.CurrentZodiacs(tGiven, loc, lang)
	resp.MoonDaysDetailed = s.moonCache.MoonDetailed(tGiven, loc, lang)

	return c.JSON(resp)
}

/*    MOON TABLE    */

func (s *Server) moonTableCurrentV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := moon.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tGiven := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableYearV1(c *fiber.Ctx) error {
	utc := c.Query("utc", "UTC:+0")
	loc, _ := moon.SetTimezoneLocFromString(utc)
	/*if err != nil {
		log.Println(err)
	}*/

	tNow := time.Now()
	year := strToInt(c.Query("year", strconv.Itoa(tNow.Year())), tNow.Year(), 0)

	tGiven := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	return s.moonTableV1(c, tGiven)
}

func (s *Server) moonTableV1(c *fiber.Ctx, tGiven time.Time) error {
	resp := MoonTable{}
	resp.Table = s.moonCache.GenerateMoonTable(tGiven)
	return c.JSON(resp.Table)
}
