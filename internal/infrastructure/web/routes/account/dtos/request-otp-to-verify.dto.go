package dtos

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type RequestOTPToVerifyInDTO struct {
	Email string `json:"email"`
}

type RequestOTPToVerifyOutDTO struct {
	SessionId string    `json:"sessionId"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func (rov *RequestOTPToVerifyInDTO) Validate() error {
	return validation.ValidateStruct(rov,
		validation.Field(&rov.Email, validation.Required, is.Email))
}
