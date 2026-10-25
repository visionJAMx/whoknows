package web

import (
	"github.com/visionJAMx/whoknows/src/models"
	"github.com/visionJAMx/whoknows/src/service"
)

// response bodies as defined in docs/openapi/swagger.json

type SearchResponse struct {
	Data []models.Page `json:"data"`
}

type SearchErrorResponse struct {
	Error string `json:"error"`
}

type RequestValidationError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type ValidationErrorDetail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}

type ValidationErrorResponse struct {
	Detail []ValidationErrorDetail `json:"detail"`
}

// errorTypes maps service problems to the contract's error type strings
var errorTypes = map[service.Problem]string{
	service.ProblemMissing:  "value_error.missing",
	service.ProblemLength:   "value_error.any_str.length",
	service.ProblemInvalid:  "value_error.email",
	service.ProblemMismatch: "value_error.mismatch",
	service.ProblemTaken:    "value_error.duplicate",
}
