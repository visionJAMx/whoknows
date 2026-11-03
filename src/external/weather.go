package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	forecastURL       = "https://api.open-meteo.com/v1/forecast"
	forecastLatitude  = "55.6761"
	forecastLongitude = "12.5683"
	forecastTimezone  = "Europe/Copenhagen"
	forecastLocation  = "Copenhagen"
	dailyParams       = "temperature_2m_max,temperature_2m_min,weather_code"
	forecastDays      = "7"
)

type Day struct {
	Date        string  `json:"date"`
	MaxTemp     float64 `json:"max_temp"`
	MinTemp     float64 `json:"min_temp"`
	WeatherCode int     `json:"weather_code"`
}

type Forecast struct {
	Location string `json:"location"`
	Days     []Day  `json:"days"`
}

type apiResponse struct {
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		WeatherCode    []int     `json:"weather_code"`
	} `json:"daily"`
}

// WeatherClient fetches the Copenhagen forecast from Open-Meteo and keeps the
// last one in memory until it's no longer today's forecast.
type WeatherClient struct {
	loc *time.Location

	mu     sync.Mutex
	cached *Forecast
}

func NewWeatherClient() (*WeatherClient, error) {
	loc, err := time.LoadLocation(forecastTimezone)
	if err != nil {
		return nil, fmt.Errorf("load %s location: %w", forecastTimezone, err)
	}

	return &WeatherClient{loc: loc}, nil
}

// Forecast returns the cached forecast if its first day is still today,
// otherwise it fetches a fresh one and caches that instead.
func (c *WeatherClient) Forecast(ctx context.Context) (*Forecast, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cached != nil && !c.stale(c.cached) {
		return c.cached, nil
	}

	forecast, err := c.fetch(ctx)
	if err != nil {
		return nil, err
	}

	c.cached = forecast
	return forecast, nil
}

// stale reports whether f's first day is no longer today, in forecastTimezone.
func (c *WeatherClient) stale(f *Forecast) bool {
	if len(f.Days) == 0 {
		return true
	}

	today := time.Now().In(c.loc).Format("2006-01-02")
	return f.Days[0].Date != today
}

func (c *WeatherClient) fetch(ctx context.Context) (*Forecast, error) {
	fmt.Println("Fetching data from api...")
	query := url.Values{
		"latitude":      {forecastLatitude},
		"longitude":     {forecastLongitude},
		"daily":         {dailyParams},
		"timezone":      {forecastTimezone},
		"forecast_days": {forecastDays},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, forecastURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create weather request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode weather response: %w", err)
	}

	forecast := Forecast{
		Location: forecastLocation,
		Days:     []Day{},
	}

	for i := range apiResp.Daily.Time {
		forecast.Days = append(forecast.Days, Day{
			Date:        apiResp.Daily.Time[i],
			MaxTemp:     apiResp.Daily.TemperatureMax[i],
			MinTemp:     apiResp.Daily.TemperatureMin[i],
			WeatherCode: apiResp.Daily.WeatherCode[i],
		})
	}

	return &forecast, nil
}
