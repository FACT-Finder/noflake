package api

import (
	"net/http"

	"github.com/FACT-Finder/noflake/database"
	"github.com/labstack/echo/v4"
)

func (a *api) GetFlakyTests(ctx echo.Context) error {
	tests, err := database.GetFlakyTests(a.db)
	if err != nil {
		return err
	}

	flakes := []FlakyTest{}

	for _, test := range tests {
		totalFails := test.TotalFails
		lastFail := test.LastFail
		flakes = append(flakes, FlakyTest{Test: test.Name, TotalFails: &totalFails, LastFail: &lastFail})
	}

	ctx.JSON(http.StatusOK, flakes)
	return nil
}
