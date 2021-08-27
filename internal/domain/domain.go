package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Request struct {
	RequestID    string `schema:"request_id"`
	SearchString string `schema:"search_string"`
}

type Response struct {
	HostToOptimalRPS map[string]int `json:"host_to_optimal_rps"`
}

func (r *Request) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(
			&r.SearchString,
			validation.Length(3, 100),
		),
	)
}
