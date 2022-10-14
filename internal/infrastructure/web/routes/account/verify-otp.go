package account

import (
	"net/http"

	"github.com/labstack/echo/v5"
	pb_rest "github.com/pocketbase/pocketbase/tools/rest"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/middlewares"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/routes/account/dtos"
)

func (r *accountResource) VerifyOTP(c echo.Context) error {
	model := &dtos.VerifyOTPIn{}
	if err := c.Bind(model); err != nil {
		return err
	}

	err := model.Validate()
	if err != nil {
		return err
	}

	rs := middlewares.GetRequestScope(c)

	if err := r.service.VerifyOTP(rs, model.Email, model.SessionId, model.OTP); err != nil {
		return pb_rest.NewApiError(http.StatusInternalServerError, "OTP is expired or not existed", err)
	}

	registeringSesisonId, expiredAt, err := r.service.CreateRegisteringSessionId(rs, model.Email)
	if err != nil {
		return pb_rest.NewApiError(http.StatusInternalServerError, "Failed to create registering session", err)
	}

	return c.JSON(http.StatusOK, &dtos.VerifyOTPOut{
		VerifiedSessionId: registeringSesisonId,
		ExpiredAt:         expiredAt,
	})
}
