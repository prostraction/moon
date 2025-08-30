package moon

import (
	"math"
	"time"
)

func GetMoonIllumination(tGiven time.Time) float64 {
	return getIlluminatedFractionOfMoon(ToJulianDate(tGiven))
}

func GetMoonPhase(before, current, after float64) (string, string) {
	switch {
	case current > 0.05 && current < 0.45 && current < after:
		return phases[0], phasesEmoji[0]
	case current >= 0.45 && current <= 0.55 && current < after:
		return phases[1], phasesEmoji[1]
	case current > 0.55 && current < 0.95 && current > before:
		return phases[2], phasesEmoji[2]
	case current >= 0.95:
		return phases[3], phasesEmoji[3]
	case current < 0.95 && current > 0.55 && current < before:
		return phases[4], phasesEmoji[4]
	case current <= 0.55 && current >= 0.45 && current < before:
		return phases[5], phasesEmoji[5]
	case current < 0.45 && current > 0.05 && current < before:
		return phases[6], phasesEmoji[6]
	case current <= 0.05:
		return phases[7], phasesEmoji[7]
	}
	return "", ""
}

func getIlluminatedFractionOfMoon(jd float64) float64 {
	const toRad = math.Pi / 180.
	T := (jd - 2451545.) / 36525.

	D := constrain(297.8501921+445267.1114034*T-0.0018819*T*T+1./545868.0*T*T*T-1./113065000.0*T*T*T*T) * toRad
	M := constrain(357.5291092+35999.0502909*T-0.0001536*T*T+1./24490000.0*T*T*T) * toRad
	Mp := constrain(134.9633964+477198.8675055*T+0.0087414*T*T+1./69699.0*T*T*T-1./14712000.0*T*T*T*T) * toRad

	i := constrain(180.-D*180./math.Pi-6.289*math.Sin(Mp)+2.1*math.Sin(M)-1.274*math.Sin(2.*D-Mp)-0.658*math.Sin(2.*D)-0.214*math.Sin(2.*Mp)-0.11*math.Sin(D)) * toRad

	k := (1. + math.Cos(i)) / 2.
	return k
}
