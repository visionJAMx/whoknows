package service

import (
	"fmt"
	"strings"
)

type Problem string

const (
	ProblemMissing  Problem = "missing"
	ProblemLength   Problem = "length"
	ProblemInvalid  Problem = "invalid"
	ProblemMismatch Problem = "mismatch"
	ProblemTaken    Problem = "taken"
)

// FieldError describes one invalid input field.
type FieldError struct {
	Field   string // name of the form field, e.g. "username"
	Problem Problem
	Message string
}

// ValidationError is returned when one or more input fields are invalid.
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	problems := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		problems[i] = fmt.Sprintf("%s: %s", f.Field, f.Problem)
	}

	return "invalid input: " + strings.Join(problems, ", ")
}
