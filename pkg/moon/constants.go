package moon

import "time"

var months = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

var phases = []string{"Waxing Crescent", "First quarter", "Waxing Gibbous", "Full Moon", "Waning Gibbous", "Third quarter", "Waning Crescent", "New Moon"}
var phasesEmoji = []string{"🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘", "🌑"}

var signs = []string{
	"Virgo", "Libra", "Scorpio", "Sagittarius",
	"Capricorn", "Aquarius", "Pisces", "Aries",
	"Taurus", "Gemini", "Cancer", "Leo",
}

var Fhour = 24.
var Fminute = 24. * 60.
var Fseconds = 24. * 60. * 60.
