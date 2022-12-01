package api

import (
	"net/http"
	"strconv"

	"github.com/FACT-Finder/noflake/database"
	"github.com/labstack/echo/v4"
)

func (a *api) flakyTests(lastNDays int) ([]FlakyTest, error) {
	tests, err := database.GetFlakyTests(a.db, lastNDays)
	if err != nil {
		return nil, err
	}

	flakes := []FlakyTest{}

	for _, test := range tests {
		test := test
		idStr := strconv.Itoa(test.ID)
		flakes = append(flakes,
			FlakyTest{
				Id:           idStr,
				Name:         test.Name,
				FailCount:    test.Fails,
				SuccessCount: test.Successes,
				Score:        test.Score,
			})
	}

	return flakes, nil
}

func (a *api) GetFlakyTests(ctx echo.Context, params GetFlakyTestsParams) error {
	lastNDays := 14
	if params.LastNDays != nil {
		lastNDays = *params.LastNDays
	}
	flakes, err := a.flakyTests(lastNDays)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, flakes)
}
