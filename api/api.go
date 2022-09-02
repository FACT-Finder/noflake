package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FACT-Finder/noflake/asset"
	"github.com/FACT-Finder/noflake/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(db *sqlx.DB, token string) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Recover())
	app.Renderer = asset.Renderer
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

	api := &api{db: db, token: token}
	wrapper := ServerInterfaceWrapper{Handler: api}

	app.POST("/report/:commit", secure(token, wrapper.AddReport))
	app.GET("/flakes", wrapper.GetFlakyTests)
	app.GET("/test/:name/fails", wrapper.GetFailedBuilds)
	app.GET("/ui", func(c echo.Context) error {
		flakes, err := api.flakyTests()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Title":  "Overview",
			"Flakes": flakes,
		})
	})
	app.GET("/ui/test/:name/fails", func(c echo.Context) error {
		name := c.Param("name")
		fails, err := api.failedBuilds(name)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "fails.html", map[string]interface{}{
			"Title": "Failed Builds " + name,
			"Fails": fails,
		})
	})
	app.StaticFS("/ui", asset.Static)

	spec, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Register(app, spec)

	return app
}

func secure(token string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		parts := strings.Fields(c.Request().Header.Get("Authorization"))

		if len(parts) == 0 {
			return unauthorized(c, "missing authorization header")
		}

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return unauthorized(c, "invalid authorization header")
		}

		if parts[1] != token {
			return unauthorized(c, "invalid token")
		}

		return handler(c)
	}
}

func unauthorized(c echo.Context, reason string) error {
	return echo.NewHTTPError(http.StatusUnauthorized, reason)
}

type api struct {
	db    *sqlx.DB
	token string
}
