package account

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	pb_rest "github.com/pocketbase/pocketbase/tools/rest"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/middlewares"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/routes/account/dtos"
)

func (r *accountResource) VerifyRegisteringMail(c echo.Context) error {
	if !r.app.Settings().Smtp.Enabled {
		return pb_rest.NewApiError(http.StatusInternalServerError, "This feature is not configured yet", errors.New("SMTP must be enabled"))
	}

	requestOTPToVerifyDTO := &dtos.RequestOTPToVerifyInDTO{}
	if err := c.Bind(requestOTPToVerifyDTO); err != nil {
		return pb_rest.NewApiError(http.StatusBadRequest, "Invalid request", err)
	}

	err := requestOTPToVerifyDTO.Validate()
	if err != nil {
		return pb_rest.NewApiError(http.StatusBadRequest, "Invalid request", err)
	}

	otpSession, expiredAt, err := r.service.RequestOtpToVerify(middlewares.GetRequestScope(c), requestOTPToVerifyDTO.Email)

	if err != nil {
		return pb_rest.NewApiError(http.StatusInternalServerError, err.Error(), nil)
	}

	c.JSON(http.StatusOK, &dtos.RequestOTPToVerifyOutDTO{
		SessionId: otpSession,
		ExpiredAt: expiredAt,
	})

	return nil
}
