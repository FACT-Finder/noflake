package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FACT-Finder/noflake/database"
	"github.com/labstack/echo/v4"
)

func (a *api) GetFailedBuilds(ctx echo.Context, testID string) error {
	flakes, err := a.failedBuilds(testID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, flakes)
}

func (a *api) failedBuilds(id string) (TestResults, error) {
	testID, err := strconv.Atoi(id)
	if err != nil {
		return TestResults{}, err
	}

	testName, err := database.GetTestName(a.db, testID)
	if err != nil {
		return TestResults{}, err
	}

	failures, err := database.GetFailures(a.db, testID)
	if err != nil {
		return TestResults{}, err
	}

	results := []TestResult{}

	for _, fail := range failures {
		uploadID := strconv.Itoa(fail.UploadID)
		commitSHA := fail.CommitSHA
		url := fail.URL
		lastFailStr := fail.Date.UTC().Format(time.RFC3339)
		results = append(results,
			TestResult{
				TestId:    id,
				UploadId:  uploadID,
				Name:      testName,
				CommitSHA: commitSHA,
				Date:      lastFailStr,
				Success:   false,
				Url:       url,
			})
	}

	return TestResults{
		TestId:  id,
		Name:    testName,
		Results: results,
	}, nil
}
