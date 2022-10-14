package middlewares

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/stewie1520/ecommerce-backend/internal/config"
	"github.com/stewie1520/ecommerce-backend/internal/tools/request"
)

type ReCaptchaModel struct {
	Token string `json:"recaptchaToken"`
}

type ReCaptchaResponse struct {
	Success bool `json:"success"`
}

const (
	GOOGLE_RECAPTCHA_ENDPOINT = "https://www.google.com/recaptcha/api/siteverify"
)

func RequireCaptcha(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		model := &ReCaptchaModel{}
		err := request.SafeBind(c, model)

		if err != nil {
			return err
		}

		if model.Token == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Token is required")
		}

		if config.AppConfig.IsDev() {
			fmt.Println("Skip captcha in dev mode")
			return next(c)
		}

		client := resty.New()
		resp, err := client.R().
			SetQueryParams(map[string]string{
				"secret":   config.AppConfig.ReCAPTCHASecretKey,
				"response": model.Token,
			}).
			SetResult(&ReCaptchaResponse{}).
			Post(GOOGLE_RECAPTCHA_ENDPOINT)

		if err != nil {
			return err
		}

		verificationErr := echo.NewHTTPError(http.StatusBadRequest, "reCAPTCHA verification failed")

		if resp.RawResponse.StatusCode != http.StatusOK {
			return verificationErr
		}

		data := resp.Result().(*ReCaptchaResponse)

		if !data.Success {
			return verificationErr
		}

		return next(c)
	}
}
