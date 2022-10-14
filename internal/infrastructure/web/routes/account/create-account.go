package account

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/middlewares"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/routes/account/dtos"
)

func (r *accountResource) CreateAccount(c echo.Context) error {
	model := &dtos.CreateAccountInDTO{}
	if err := c.Bind(model); err != nil {
		return err
	}

	err := model.Validate()
	if err != nil {
		return err
	}

	rs := middlewares.GetRequestScope(c)
	err = r.service.VerifyRegisteringSessionId(rs, model.VerifiedSessionId, model.Email)
	if err != nil {
		return err
	}

	userId, err := r.service.CreateUserByEmail(
		rs,
		model.Email,
		model.Password,
		model.Name,
		model.Birthday,
		true,
	)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, &dtos.CreateAccountOutDTO{
		UserId: userId,
		Name:   model.Name,
	})

	return nil
}
