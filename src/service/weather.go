package service

import (
	"context"

	"github.com/visionJAMx/whoknows/src/external"
)

type WeatherClient interface {
	Forecast(ctx context.Context) (*external.Forecast, error)
}

type WeatherService struct {
	client WeatherClient
}

func NewWeatherService(client WeatherClient) *WeatherService {
	return &WeatherService{client: client}
}

func (s *WeatherService) Forecast(ctx context.Context) (*external.Forecast, error) {
	return s.client.Forecast(ctx)
}
