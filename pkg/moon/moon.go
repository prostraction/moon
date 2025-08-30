package moon

import (
	"time"
)

func CalcMoonNumber(yearGiven int) int {
	d := 6 // 2000
	mod := 19
	yearCurrent := 2000

	if yearCurrent < yearGiven {
		for yearCurrent+mod < yearGiven {
			yearCurrent += mod
		}
		for yearCurrent < yearGiven {
			yearCurrent++
			d++
			if d > 19 {
				d = 1
			}
		}
	} else {
		for yearCurrent-mod > yearGiven {
			yearCurrent -= mod
		}
		for yearCurrent > yearGiven {
			yearCurrent--
			d--
			if d < 1 {
				d = 19
			}
		}
	}
	return d
}

func Gen(tGiven time.Time) ([]*MoonTableElement, time.Duration, float64, string) {
	//

	moonTable := CreateMoonTable(tGiven)
	moonDays := GetMoonDays(tGiven, moonTable)
	moonIllumination := GetMoonIllumination(tGiven)

	//zodiacPosition := int((jdIllumination*360)/30) % 12

	return moonTable, moonDays, moonIllumination, "not working" //getZodiacSign(zodiacPosition)
}

func getZodiacSign(position int) string {
	signs := []string{
		"Aries", "Taurus", "Gemini", "Cancer",
		"Leo", "Virgo", "Libra", "Scorpio",
		"Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	return signs[position]
}
