package moon

import "time"

type MoonTableElement struct {
	NewMoon      time.Time
	FirstQuarter time.Time
	FullMoon     time.Time
	LastQuarter  time.Time
	t1           float64
	t2           float64
}

type Cache struct {
	tables map[string][]*MoonTableElement
}

type PhaseResp struct {
	Name          string
	NameLocalized string
	Emoji         string
}

type Zodiac struct {
	Name          string
	NameLocalized string
	Emoji         string
}

type ZodiacDetailed struct {
	Name          string
	NameLocalized string
	Emoji         string
	Begin         time.Time
	End           time.Time
}

type Zodiacs struct {
	Count  int
	Zodiac []ZodiacDetailed
}

type MoonDay struct {
	Begin time.Time
	End   time.Time
}

type MoonDaysDetailed struct {
	Count int
	Day   []MoonDay
}
