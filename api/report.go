package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FACT-Finder/noflake/database"
	"github.com/FACT-Finder/noflake/model"
	"github.com/joshdk/go-junit"
	"github.com/labstack/echo/v4"
)

func (a *api) AddReport(ctx echo.Context, commitSha string, params AddReportParams) error {
	commit, err := database.CreateOrGetCommit(a.db, model.Commit{CommitSha: commitSha})
	if err != nil {
		return err
	}

	upload, err := database.CreateUpload(a.db, model.Upload{CommitID: *commit.ID, Time: time.Now().UTC(), URL: params.Url})
	if err != nil {
		return err
	}

	r := ctx.Request()
	err = r.ParseMultipartForm(200000)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid multipart request: %s", err))
	}

	formdata := r.MultipartForm

	files := formdata.File["report"]

	tests := make([]model.TestResult, 0)
	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("couldn't open file: %s", err))
		}
		defer file.Close()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("couldn't read file: %s", err))
		}

		suites, err := junit.IngestReader(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("couldn't parse file %s: %s", files[i].Filename, err))
		}

		for _, suite := range suites {
			for _, testcase := range suite.Tests {
				if testcase.Status == junit.StatusSkipped {
					continue
				}
				test := model.TestResult{
					Name:    fmt.Sprintf("%s.%s", suite.Name, testcase.Name),
					Success: testcase.Status == junit.StatusPassed,
				}
				if testcase.Error != nil {
					output := testcase.Error.Error()
					test.Output = &output
				}
				tests = append(tests, test)
			}
		}
	}

	err = database.InsertTests(a.db, tests, *upload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("couldn't store test results: %s", err))
	}

	return ctx.NoContent(http.StatusNoContent)
}
