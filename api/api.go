package api

import (
	"fmt"
	"net/http"

	"github.com/FACT-Finder/noflake/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(db *sqlx.DB) *echo.Echo {
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

	wrapper := ServerInterfaceWrapper{Handler: &api{db: db}}

	app.POST("/report/:commit", wrapper.AddReport)
	app.GET("/flakes", wrapper.GetFlakyTests)

	spec, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Register(app, spec)

	return app
}

type api struct {
	db *sqlx.DB
}
