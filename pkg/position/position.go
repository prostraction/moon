package position

import (
	"encoding/json"
	"fmt"
	"io"
	jt "moon/pkg/julian-time"
	"net/http"
	"net/url"
	"time"
)
/*
type Position struct {
	client  *http.Client
	baseURL string
}

func New() *Position {
	return &Position{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: "http://localhost:9997/",
	}
}

func (p *Position) WithHTTPClient(client *http.Client) *Position {
	p.client = client
	return p
}
func (p *Position) WithBaseURL(baseURL string) *Position {
	p.baseURL = baseURL
	return p
}*/

// DayResponse
type DayResponse struct {
	Status     string     `json:"status"`
	Parameters Parameters `json:"parameters"`
	Data       DayData    `json:"data"`
	Range      string     `json:"range"`
}

// MonthResponse
type MonthResponse struct {
	Status     string     `json:"status"`
	Parameters Parameters `json:"parameters"`
	Data       []DayData  `json:"data"`
	Range      string     `json:"range"`
	DaysCount  int        `json:"days_count"`
}

// input
type Parameters struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  int     `json:"timezone"`
	UTCOffset string  `json:"utc_offset"`
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	Day       int     `json:"day,omitempty"`
}

// resp for 1 daya
type DayData struct {
	Date       string  `json:"date"`
	Moonrise   int64   `json:"moonrise"`
	Moonset    int64   `json:"moonset"`
	Meridian   int64   `json:"meridian"`
	IsMoonRise bool    `json:"isMoonRise"`
	IsMoonSet  bool    `json:"isMoonSet"`
	IsMeridian bool    `json:"isMeridian"`
	DistanceKm float64 `json:"distance_km"`
}

func GetRisesMonthly(year, month int, loc *time.Location, location ...float64) (*MonthResponse, error) {
	lat, lon := p.parseLocation(location)

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

	url := p.baseURL + "?" + params.Encode()

	resp, err := p.client.Get(url)
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

func GetRisesDay(year, month, day int, loc *time.Location, location ...float64) (*DayResponse, error) {
	lat, lon := p.parseLocation(location)

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
	params.Add("day", fmt.Sprintf("%d", day))

	url := p.baseURL + "?" + params.Encode()
	resp, err := p.client.Get(url)
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

	return &dayResponse, nil
}

func (p *Position) parseLocation(location []float64) (lat, lon float64) {
	lat = 51.08 // Astana
	lon = 71.26 // Astana

	if len(location) == 2 {
		lat = location[0]
		lon = location[1]
	}

	return lat, lon
}

func (d *DayData) MoonriseTime() time.Time {
	return time.Unix(d.Moonrise, 0)
}
func (d *DayData) MoonsetTime() time.Time {
	return time.Unix(d.Moonset, 0)
}
func (d *DayData) MeridianTime() time.Time {
	return time.Unix(d.Meridian, 0)
}
