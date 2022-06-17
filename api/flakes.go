package api

import (
	"net/http"
	"time"

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
		lastFailStr := lastFail.UTC().Format(time.RFC3339)
		flakes = append(flakes,
			FlakyTest{Test: test.Name, TotalFails: &totalFails, LastFail: &lastFailStr})
	}

	ctx.JSON(http.StatusOK, flakes)
	return nil
}
