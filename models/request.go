package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ReserveRequest struct {
	UserID string `header:"userID"`
}

func (a ReserveRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UserID, validation.Required, is.UUID),
	)
}

type SnatchRequest struct {
	UserID string `header:"userID"`
}

func (a SnatchRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UserID, validation.Required, is.UUID),
	)
}
