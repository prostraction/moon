package moon

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func ToJulianDate(t time.Time) float64 {
	m := 1
	for i := range months {
		if t.Month() == months[i] {
			m = i + 1
		}
	}
	if m < 3 {
		t = t.AddDate(-1, 0, 0)
		m += 12
	}

	A := float64(t.Year()) / 100
	B := A / 4
	C := 2 - A + B
	E := 365.25 * float64(t.Year()+4716)
	F := 30.6001 * (float64(m) + 1)
	return C + float64(float64(t.Day()+t.Hour()/24.+t.Minute()/60)) + E + F - 1524.5
}

func FromJulianDate(j float64, loc *time.Location) time.Time {
	datey, datem, dated := jyear(j)

	//j += gmtOffset / 24.0 // 5 timezones west of UTC
	//j += (30.0 / (24 * 60 * 60))

	timeh, timem, times := jhms(j)

	t := time.Date(datey, getMonth(datem), dated, timeh, timem, times, 0, time.UTC)
	t = t.In(loc)
	return t
}

func SetTimezoneLocFromString(utc string) (*time.Location, error) {
	if utc == "" {
		return time.UTC, nil
	}

	// Remove "UTC" prefix if present and convert to lowercase for case-insensitive matching
	normalized := strings.ToLower(strings.TrimPrefix(utc, "UTC"))
	normalized = strings.TrimSpace(normalized)

	// Handle cases like "UTC", "+0", "-0", "0"
	if normalized == "" || normalized == "+0" || normalized == "-0" || normalized == "0" {
		return time.UTC, nil
	}

	// Check if it starts with + or -
	sign := 1
	if strings.HasPrefix(normalized, "+") {
		sign = 1
		normalized = normalized[1:]
	} else if strings.HasPrefix(normalized, "-") {
		sign = -1
		normalized = normalized[1:]
	}

	var hours, minutes int
	var err error

	// Handle cases with colon separator (e.g., "05:30", "5:30")
	if strings.Contains(normalized, ":") {
		parts := strings.Split(normalized, ":")
		if len(parts) != 2 {
			return time.UTC, fmt.Errorf("invalid timezone format: %s", utc)
		}

		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return time.UTC, fmt.Errorf("invalid hours: %s", parts[0])
		}

		minutes, err = strconv.Atoi(parts[1])
		if err != nil || minutes < 0 || minutes >= 60 {
			return time.UTC, fmt.Errorf("invalid minutes: %s", parts[1])
		}

	} else {
		// Handle cases without colon
		switch len(normalized) {
		case 1, 2: // Just hours (e.g., "5", "12")
			hours, err = strconv.Atoi(normalized)
			if err != nil {
				return time.UTC, fmt.Errorf("invalid hours: %s", normalized)
			}
			minutes = 0

		case 3: // Hours + minutes (e.g., "530" -> 5 hours, 30 minutes)
			hours, err = strconv.Atoi(normalized[:1])
			if err != nil {
				return time.UTC, fmt.Errorf("invalid hours: %s", normalized[:1])
			}
			minutes, err = strconv.Atoi(normalized[1:])
			if err != nil || minutes < 0 || minutes >= 60 {
				return time.UTC, fmt.Errorf("invalid minutes: %s", normalized[1:])
			}

		case 4: // Hours + minutes (e.g., "0530" -> 5 hours, 30 minutes)
			hours, err = strconv.Atoi(normalized[:2])
			if err != nil {
				return time.UTC, fmt.Errorf("invalid hours: %s", normalized[:2])
			}
			minutes, err = strconv.Atoi(normalized[2:])
			if err != nil || minutes < 0 || minutes >= 60 {
				return time.UTC, fmt.Errorf("invalid minutes: %s", normalized[2:])
			}

		default:
			return time.UTC, fmt.Errorf("invalid timezone format: %s", utc)
		}
	}

	// Validate hours range
	if hours < 0 || hours > 23 {
		return time.UTC, fmt.Errorf("hours out of range (0-23): %d", hours)
	}

	// Calculate total seconds offset
	totalSeconds := sign * (hours*3600 + minutes*60)

	// Create location name
	locationName := fmt.Sprintf("UTC%s%d:%02d", getSignPrefix(sign), hours, minutes)
	if minutes == 0 {
		locationName = fmt.Sprintf("UTC%s%d", getSignPrefix(sign), hours)
	}

	return time.FixedZone(locationName, totalSeconds), nil
}

// JYMD - Convert Julian time to year, months, and days
func jyear(td float64) (int, int, int) {
	td += 0.5
	z := math.Floor(td)
	f := td - z

	var a float64
	if z < 2299161.0 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}

	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	mm := int(math.Floor(e))
	if mm >= 14 {
		mm -= 13
	} else {
		mm -= 1
	}

	year := int(math.Floor(c))
	if mm > 2 {
		year -= 4716
	} else {
		year -= 4715
	}

	day := int(math.Floor(b - d - math.Floor(30.6001*e) + f))

	return year, mm, day
}

// JHMS - Convert Julian time to hour, minutes, and seconds
func jhms(j float64) (int, int, int) {
	j += 0.5 // Astronomical to civil
	ij := (j - math.Floor(j)) * 86400.0
	hours := math.Floor(ij / 3600)
	minutes := math.Floor((ij / 60))
	seconds := math.Floor(ij)
	return int(hours), int(math.Mod(minutes, 60)), int(math.Mod(seconds, 60))
}

func getMonth(datem int) time.Month {
	datem = min(max(datem-1, 0), 11)
	return months[datem]
}
