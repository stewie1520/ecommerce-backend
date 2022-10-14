package dtos

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type VerifyOTPIn struct {
	Email     string `json:"email"`
	SessionId string `json:"sessionId"`
	OTP       string `json:"otp"`
}

type VerifyOTPOut struct {
	VerifiedSessionId string    `json:"verifiedSessionId"`
	ExpiredAt         time.Time `json:"expiredAt"`
}

func (vo *VerifyOTPIn) Validate() error {
	return validation.ValidateStruct(vo,
		validation.Field(&vo.Email, validation.Required, is.Email),
		validation.Field(&vo.SessionId, validation.Required),
		validation.Field(&vo.OTP, validation.Required))
}
