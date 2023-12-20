package pkg

import "github.com/labstack/echo/v4"

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (svc *AdminSvc) AddSVCToEchoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("svc", svc)
		return next(c)
	}
}
