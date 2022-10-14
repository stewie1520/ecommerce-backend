package account

import (
	"github.com/stewie1520/ecommerce-backend/internal/core"
)

func (as *AccountService) VerifyOTP(rs core.RequestScope, email string, sessionId string, otp string) error {
	val, err := as.cache.Get(rs.GetContext(), sessionId).Result()
	if err != nil {
		return err
	}

	otpValue := as.formatOTP(email, otp)
	err = as.hash.CompareHashAndPassword(val, otpValue)
	if err != nil {
		return err
	}

	return nil
}
