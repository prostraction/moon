package moon

import (
	"math"
	"time"
)

var months = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

var phases = []string{"Waxing Crescent", "First quarter", "Waxing Gibbous", "Full Moon", "Waning Gibbous", "Third quarter", "Waning Crescent", "New Moon"}
var phasesEmoji = []string{"🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘", "🌑"}

var signs = []string{
	"Virgo", "Libra", "Scorpio", "Sagittarius",
	"Capricorn", "Aquarius", "Pisces", "Aries",
	"Taurus", "Gemini", "Cancer", "Leo",
}

const Fhour = 24.
const Fminute = 24. * 60.
const Fseconds = 24. * 60. * 60.

const toRad = math.Pi / 180.

var tables map[time.Time][]*MoonTableElement
