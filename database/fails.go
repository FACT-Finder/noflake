package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type TestFailure struct {
	Date      time.Time `db:"last_fail"`
	Output    *string   `db:"output"`
	CommitSHA string    `db:"commit_sha"`
	URL       *string   `db:"url"`
}

func GetFailures(db *sqlx.DB, name string) ([]TestFailure, error) {
	rows, err := db.Queryx(`
	SELECT
		commits.commit_sha,
		uploads.url,
		results.output,
		uploads.time
	from results
	LEFT JOIN commits on commits.id = results.commit_id
	LEFT JOIN uploads on uploads.id = results.upload_id
	LEFT JOIN tests on tests.id = results.test_id
	WHERE
		tests.name = ? and results.success = 0
	ORDER BY uploads.time desc
	`, name)
	if err != nil {
		return nil, err
	}

	failures := []TestFailure{}
	for rows.Next() {
		var failure TestFailure
		var dateTimestamp int64
		err = rows.Scan(&failure.CommitSHA, &failure.URL, &failure.Output, &dateTimestamp)
		if err != nil {
			return nil, err
		}
		failure.Date = time.Unix(dateTimestamp, 0)
		failures = append(failures, failure)
	}

	return failures, nil
}
