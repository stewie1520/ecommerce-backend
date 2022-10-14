package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	pb_core "github.com/pocketbase/pocketbase/core"
	service_account "github.com/stewie1520/ecommerce-backend/internal/core/services/account"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/middlewares"
	route_account "github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/routes/account"
	"go.uber.org/zap"
)

type PocketBaseCustomRoute struct {
	app *pocketbase.PocketBase
}

func RegisterRoutes(app *pocketbase.PocketBase, accountService service_account.IAccountService) {
	app.OnBeforeServe().Add(func(e *pb_core.ServeEvent) error {
		e.Router.Use(middlewares.WithRequestScope(app))
		logger, _ := zap.NewProduction()
		e.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
				)

				return nil
			},
		}))

		apiGroup := e.Router.Group("/api")

		route_account.ServeAccountResource(apiGroup.Group("/account"), app, accountService)
		return nil
	})

}
