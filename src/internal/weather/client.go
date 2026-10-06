package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIResponse struct {
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		WeatherCode    []int     `json:"weather_code"`
	} `json:"daily"`
}

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

func GetForecast(ctx context.Context) (*Forecast, error) {
	url := "https://api.open-meteo.com/v1/forecast?latitude=55.6761&longitude=12.5683&daily=temperature_2m_max,temperature_2m_min,weather_code&timezone=Europe%2FCopenhagen&forecast_days=7"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var apiResponse APIResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("decode weather response: %w", err)
	}

	forecast := Forecast{
		Location: "Copenhagen",
		Days:     []Day{},
	}

	for i := range apiResponse.Daily.Time {
		forecast.Days = append(forecast.Days, Day{
			Date:        apiResponse.Daily.Time[i],
			MaxTemp:     apiResponse.Daily.TemperatureMax[i],
			MinTemp:     apiResponse.Daily.TemperatureMin[i],
			WeatherCode: apiResponse.Daily.WeatherCode[i],
		})
	}

	return &forecast, nil
}
