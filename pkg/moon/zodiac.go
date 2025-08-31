package moon

import (
	"time"
)

func (c *Cache) CurrentZodiacs(tGiven time.Time, loc *time.Location) (*Zodiacs, Zodiac, Zodiac, Zodiac) {
	zods := new(Zodiacs)

	zodiacBegin := Zodiac{}
	zodiacCurrent := Zodiac{}
	zodiacEnd := Zodiac{}

	dayBeginTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	dayEndTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()+1, 0, 0, 0, 0, loc)

	moonTable := c.CreateMoonTable(tGiven)
	beginMoonDays := GetMoonDays(dayBeginTime, moonTable)
	currentMoonDays := GetMoonDays(tGiven, moonTable)
	endMoonDays := GetMoonDays(dayEndTime, moonTable)

	zodiacPositionBegin := int((beginMoonDays.Minutes()/Fminute*360.)/30.) / 30. % 12
	zodiacPositionCurrent := int((currentMoonDays.Minutes()/Fminute*360.)/30.) / 30. % 12
	zodiacPositionEnd := int((endMoonDays.Minutes()/Fminute*360.)/30.) / 30. % 12

	zodiacBegin.Name, zodiacBegin.Emoji = getZodiacResp(zodiacPositionBegin)
	zodiacCurrent.Name, zodiacCurrent.Emoji = getZodiacResp(zodiacPositionCurrent)
	zodiacEnd.Name, zodiacEnd.Emoji = getZodiacResp(zodiacPositionEnd)

	return zods, zodiacBegin, zodiacCurrent, zodiacEnd
}

func getZodiacResp(position int) (string, string) {
	if position >= 0 && position < len(signs) && position < len(signsEmoji) {
		return signs[position], signsEmoji[position]
	}
	return "", ""
}
