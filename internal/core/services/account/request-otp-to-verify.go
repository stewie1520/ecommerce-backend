package account

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/stewie1520/ecommerce-backend/internal/config"
	"github.com/stewie1520/ecommerce-backend/internal/core"
	"github.com/stewie1520/ecommerce-backend/internal/templates"
)

func (as *AccountService) RequestOtpToVerify(rs core.RequestScope, email string) (string, time.Time, error) {
	pbApp := rs.PbApp()
	if !pbApp.Settings().EmailAuth.Enabled {
		return "", time.Time{}, errors.New("email/password authentication is not enabled")
	}

	existedUser, err := pbApp.Dao().FindUserByEmail(email)
	if existedUser != nil {
		return "", time.Time{}, errors.New("email already exists")
	}

	otpVal, otp := as.generateOtpValue(email)
	h, err := as.hash.HashPassword(otpVal)

	if err != nil {
		return "", time.Time{}, err
	}

	hashedOtp := string(h)
	aliveDuration := time.Minute * 5
	expiredAt := time.Now().Add(aliveDuration)

	otpKey := as.generateOtpKey()
	err = as.cache.Set(rs.GetContext(), otpKey, hashedOtp, aliveDuration).Err()
	if err != nil {
		return "", time.Time{}, err
	}

	html, err := templates.ResolveVerifyOTPMailTemplate(otp, pbApp.Settings().Meta.AppName)
	if err != nil {
		return "", time.Time{}, err
	}

	if config.AppConfig.IsDev() {
		fmt.Printf("Skip sending email in dev mode, OTP: %s, email: %s\n", otp, email)
		return otpKey, expiredAt, nil
	}

	// Send mail
	mailClient := pbApp.NewMailClient()
	err = mailClient.Send(
		mail.Address{
			Name:    pbApp.Settings().Meta.SenderName,
			Address: pbApp.Settings().Meta.SenderAddress,
		},
		mail.Address{Address: email},
		"Verify your email",
		html,
		nil,
	)

	if err != nil {
		return "", time.Time{}, err
	}

	return otpKey, expiredAt, err
}
