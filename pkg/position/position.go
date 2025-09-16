package position

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	jt "moon/pkg/julian-time"
	"net/http"
	"net/url"
	"time"
)

// PositionResponse
type PositionResponse struct {
	MoonEvent
	DistanceKm float64 `json:"DistanceKm"`
}

// DayResponse
type DayResponse struct {
	Status     string     `json:"Status"`
	Parameters Parameters `json:"Parameters"`
	Data       *DayData   `json:"Data"`
	Range      string     `json:"Range"`
}

// MonthResponse
type MonthResponse struct {
	Status     string     `json:"Status"`
	Parameters Parameters `json:"Parameters"`
	Data       []DayData  `json:"Data"`
	Range      string     `json:"Range"`
	DaysCount  int        `json:"DaysCount"`
}

// input
type Parameters struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Timezone  int     `json:"Timezone"`
	UTCOffset string  `json:"UtcOffset"`
	Year      int     `json:"Year"`
	Month     int     `json:"Month"`
	Day       int     `json:"Day,omitempty"`
}

// resp for 1 day
type MoonEvent struct {
	Timestamp       int64     `json:"Timestamp"`
	TimeISO         time.Time `json:"TimeISO,omitempty"`
	AzimuthDegrees  float64   `json:"AzimuthDegrees"`
	AltitudeDegrees float64   `json:"AltitudeDegrees"`
	Direction       string    `json:"Direction"`
}

type DayData struct {
	Moonrise   *MoonEvent `json:"Moonrise,omitempty"`
	Moonset    *MoonEvent `json:"Moonset,omitempty"`
	Meridian   *MoonEvent `json:"Meridian,omitempty"`
	DistanceKm float64    `json:"DistanceKm"`
	IsMoonRise bool       `json:"IsMoonRise"`
	IsMoonSet  bool       `json:"IsMoonSet"`
	IsMeridian bool       `json:"IsMeridian"`
}

func GetRisesMonthly(year, month int, loc *time.Location, precision int, location ...float64) (*MonthResponse, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	h := 0
	if loc != nil {
		jth, _, err := jt.GetTimeFromLocation(loc)
		if err == nil {
			h = jth
		}
	}

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	params.Add("utc", fmt.Sprintf("%d", h))
	params.Add("year", fmt.Sprintf("%d", year))
	params.Add("month", fmt.Sprintf("%d", month))
	params.Add("precision", fmt.Sprintf("%d", precision))

	url := baseURL + "?" + params.Encode()
	client := &http.Client{Timeout: 69 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var monthResponse MonthResponse
	if err := json.Unmarshal(body, &monthResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &monthResponse, nil
}

func GetRisesDay(year, month, day int, loc *time.Location, precision int, location ...float64) (*DayData, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	h, m := 0, 0
	if loc != nil {
		jth, jtm, err := jt.GetTimeFromLocation(loc)
		if err == nil {
			h = jth
			m = jtm
		}
	}

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	params.Add("utc_hours", fmt.Sprintf("%d", h))
	params.Add("utc_minutes", fmt.Sprintf("%d", m))
	params.Add("year", fmt.Sprintf("%d", year))
	params.Add("month", fmt.Sprintf("%d", month))
	params.Add("day", fmt.Sprintf("%d", day))
	params.Add("precision", fmt.Sprintf("%d", precision))

	url := baseURL + "daily" + "?" + params.Encode()
	client := &http.Client{Timeout: 69 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var dayResponse DayResponse
	if err := json.Unmarshal(body, &dayResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if dayResponse.Data.Meridian != nil {
		timestampToGoTime(dayResponse.Data.Meridian, loc)
		dayResponse.Data.Meridian.Timestamp = dayResponse.Data.Meridian.TimeISO.Unix()
	}
	if dayResponse.Data.Moonrise != nil {
		timestampToGoTime(dayResponse.Data.Moonrise, loc)
		dayResponse.Data.Moonrise.Timestamp = dayResponse.Data.Moonrise.TimeISO.Unix()
	}
	if dayResponse.Data.Moonset != nil {
		timestampToGoTime(dayResponse.Data.Moonset, loc)
		dayResponse.Data.Moonset.Timestamp = dayResponse.Data.Moonset.TimeISO.Unix()
	}

	return dayResponse.Data, nil
}

func GetMoonPosition(tGiven time.Time, loc *time.Location, precision int, location ...float64) (*PositionResponse, error) {
	lat, lon, err := parseLocation(location)
	if err != nil {
		return nil, err
	}

	h, m := 0, 0
	if loc != nil {
		tGiven = tGiven.In(loc)
		jth, jtm, err := jt.GetTimeFromLocation(loc)
		if err == nil {
			h = jth
			m = jtm
		}
	}

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	params.Add("utc_hours", fmt.Sprintf("%d", h))
	params.Add("utc_minutes", fmt.Sprintf("%d", m))
	params.Add("year", fmt.Sprintf("%d", tGiven.Year()))
	params.Add("month", fmt.Sprintf("%d", int(tGiven.Month())))
	params.Add("day", fmt.Sprintf("%d", tGiven.Day()))
	params.Add("hour", fmt.Sprintf("%d", tGiven.Hour()))
	params.Add("minute", fmt.Sprintf("%d", tGiven.Minute()))
	params.Add("second", fmt.Sprintf("%d", tGiven.Second()))
	params.Add("precision", fmt.Sprintf("%d", precision))

	url := baseURL + "position" + "?" + params.Encode()
	client := &http.Client{Timeout: 69 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var pos *PositionResponse
	if err := json.Unmarshal(body, &pos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if pos != nil {
		timestampToGoTime(&pos.MoonEvent, loc)
	}

	return pos, nil
}

func parseLocation(location []float64) (lat, lon float64, err error) {
	if len(location) == 2 {
		lat = location[1]
		lon = location[0]
	} else {
		return 0, 0, errors.New("no location prodived")
	}

	return lat, lon, nil
}

func timestampToGoTime(ev *MoonEvent, loc *time.Location) {
	utcTime := time.Unix(ev.Timestamp, 0).UTC()
	ev.TimeISO = utcTime
	if loc != nil {
		ev.TimeISO = ev.TimeISO.In(loc)
	}
}
