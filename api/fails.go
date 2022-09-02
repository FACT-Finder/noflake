package api

import (
	"net/http"
	"time"

	"github.com/FACT-Finder/noflake/database"
	"github.com/labstack/echo/v4"
)

func (a *api) GetFailedBuilds(ctx echo.Context, name string) error {
	flakes, err := a.failedBuilds(name)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, flakes)
}

func (a *api) failedBuilds(name string) ([]FailedBuild, error) {
	failures, err := database.GetFailures(a.db, name)
	if err != nil {
		return nil, err
	}

	result := []FailedBuild{}

	for _, fail := range failures {
		commitSHA := fail.CommitSHA
		output := fail.Output
		url := fail.URL
		lastFailStr := fail.Date.UTC().Format(time.RFC3339)
		result = append(result,
			FailedBuild{CommitSHA: commitSHA, Date: lastFailStr, Output: output, Url: url})
	}

	return result, nil
}
