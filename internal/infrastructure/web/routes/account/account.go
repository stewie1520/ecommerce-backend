package account

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	pb_apis "github.com/pocketbase/pocketbase/apis"
	service_account "github.com/stewie1520/ecommerce-backend/internal/core/services/account"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/middlewares"
)

type accountResource struct {
	service service_account.IAccountService

	app *pocketbase.PocketBase
}

func ServeAccountResource(route *echo.Group, app *pocketbase.PocketBase, service service_account.IAccountService) error {
	r := &accountResource{
		service: service,
		app:     app,
	}

	authRoute := route.Group("/auth", pb_apis.RequireGuestOnly())
	authRoute.POST("/send-otp", r.VerifyRegisteringMail, middlewares.RequireCaptcha)
	authRoute.POST("/verify-otp", r.VerifyOTP)
	authRoute.POST("/create-user", r.CreateAccount)
	return nil
}
