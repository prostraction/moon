package server

import (
	"moon/pkg/moon"
	"time"
)

type MoonTable struct {
	Table []*moon.MoonTableElement
}

type MoonStat struct {
	MoonDays     float64
	Illumination float64
	Phase        moon.PhaseResp
	Zodiac       moon.Zodiac
}

type FullInfo struct {
	MoonDaysBegin   float64
	MoonDaysEnd     float64
	MoonDaysCurrent float64

	IlluminationBeginDay float64
	IlluminationCurrent  float64
	IlluminationEndDay   float64
}

type MoonDay struct {
	Begin time.Time
	End   time.Time
}

type MoonPhaseResponse struct {
	BeginDay     *MoonStat
	CurrentState *MoonStat
	EndDay       *MoonStat

	MoonDaysDetailed *moon.MoonDaysDetailed
	ZodiacDetailed   *moon.Zodiacs

	info *FullInfo
}
