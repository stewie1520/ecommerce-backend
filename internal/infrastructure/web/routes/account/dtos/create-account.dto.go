package dtos

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type CreateAccountInDTO struct {
	VerifiedSessionId string `json:"verifiedSessionId"`
	Email             string `json:"email"`
	Birthday          string `json:"birthday"`
	Password          string `json:"password"`
	Name              string `json:"name"`
}

type CreateAccountOutDTO struct {
	UserId string `json:"user"`
	Name   string `json:"name"`
}

func (ca *CreateAccountInDTO) Validate() error {
	return validation.ValidateStruct(ca,
		validation.Field(&ca.VerifiedSessionId, validation.Required),
		validation.Field(&ca.Email, validation.Required, is.Email),
		validation.Field(&ca.Birthday, validation.Required, validation.Date(time.RFC3339)),
		validation.Field(&ca.Password, validation.Required, validation.Min(8)),
		validation.Field(&ca.Name, validation.Required))
}
