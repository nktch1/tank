package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	RequestID    string
	SearchString string
}

type Response struct {
	HostToOptimalRPS map[string]int
}

func (r *Request) Validate() error {
	if r.RequestID == "" {
		r.RequestID = uuid.NewV4().String()
	}

	return validation.ValidateStruct(
		r,

		validation.Field(
			&r.SearchString,
			validation.Length(6, 100),
		),
	)
}
