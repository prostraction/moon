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
	Name  string
	Emoji string
}

type Zodiac struct {
	Name  string
	Emoji string
}

type ZodiacDetailed struct {
	Name  string
	Emoji string
	Begin time.Time
	End   time.Time
}

type Zodiacs struct {
	Count  int
	Zodiac []ZodiacDetailed
}
