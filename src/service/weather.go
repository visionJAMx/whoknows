package service

import "context"

type WeatherService struct {
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

func (s *WeatherService) WeatherPage(ctx context.Context) (string, error) {

	// TODO: call weather client

	return "GET /weather", nil
}

func (s *WeatherService) Forecast(ctx context.Context) (string, error) {

	// TODO: call weather client

	return "GET /api/weather", nil
}
