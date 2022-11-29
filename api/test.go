package api

import (
	"net/http"
	"strconv"

	"github.com/FACT-Finder/noflake/database"
	"github.com/labstack/echo/v4"
)

func (a *api) GetTestResult(ctx echo.Context, testID, uploadID string) error {
	output, err := a.testResult(testID, uploadID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, output)
}

func (a *api) testResult(testIDStr, uploadIDStr string) (TestResult, error) {
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		return TestResult{}, err
	}

	uploadID, err := strconv.Atoi(uploadIDStr)
	if err != nil {
		return TestResult{}, err
	}

	testResult, err := database.GetTestResult(a.db, testID, uploadID)
	if err != nil {
		return TestResult{}, err
	}

	return TestResult{
		TestId:    testIDStr,
		UploadId:  uploadIDStr,
		Name:      testResult.Name,
		CommitSHA: testResult.CommitSHA,
		Url:       testResult.URL,
		Success:   testResult.Success,
		Output:    testResult.Output,
	}, nil
}
