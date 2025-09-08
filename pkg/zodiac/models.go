package zodiac

import "time"

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

type Zodiac struct {
	Name          string
	NameLocalized string
	Emoji         string
}
