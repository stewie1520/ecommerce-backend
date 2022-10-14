package middlewares

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/stewie1520/ecommerce-backend/internal/core"
)

const requestScopeContextKey = "requestScopeContextKey"

func WithRequestScope(app *pocketbase.PocketBase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			// X-Request-Id Header is propagated by proxy server e.g nginx
			requestID := c.Request().Header.Get("X-Request-Id")

			rs := core.NewRequestScope(requestID, ctx, app)

			c.Set(requestScopeContextKey, rs)
			return next(c)
		}
	}
}

func GetRequestScope(c echo.Context) core.RequestScope {
	return c.Get(requestScopeContextKey).(core.RequestScope)
}
