package api

import (
	"fmt"
	"net/http"

	"github.com/FACT-Finder/noflake/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(db *sqlx.DB, token string) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Recover())
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := ""
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprint(he.Message)
		} else {
			message = err.Error()
		}
		c.JSON(code, &ApiError{
			Error:       http.StatusText(code),
			Description: message,
		})
	}

	wrapper := ServerInterfaceWrapper{
		Handler: &api{
			db:    db,
			token: token,
		},
	}

	app.POST("/report/:commit", secure(token, wrapper.AddReport))
	app.GET("/flakes", wrapper.GetFlakyTests)

	spec, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Register(app, spec)

	return app
}


func secure(token string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, pass, ok := c.Request().BasicAuth()

		if !ok {
			return unauthorized(c, "noflake", "no credentials provided")
		}

		if user != "token" || pass != token{
			return unauthorized(c, "noflake", "wrong credentials")
		}

		return handler(c)
	}
}

func unauthorized(c echo.Context, realm, reason string) error {
	c.Response().Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	return echo.NewHTTPError(http.StatusUnauthorized, reason)
}

type api struct {
	db    *sqlx.DB
	token string
}
