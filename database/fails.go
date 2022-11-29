package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type TestFailure struct {
	UploadID  int       `db:"upload_id"`
	Date      time.Time `db:"last_fail"`
	CommitSHA string    `db:"commit_sha"`
	URL       *string   `db:"url"`
}

func GetFailures(db *sqlx.DB, testID int) ([]TestFailure, error) {
	rows, err := db.Queryx(`
	SELECT
		uploads.id as upload_id,
		commits.commit_sha,
		uploads.url,
		uploads.time
	from results
	LEFT JOIN commits on commits.id = results.commit_id
	LEFT JOIN uploads on uploads.id = results.upload_id
	LEFT JOIN tests on tests.id = results.test_id
	WHERE
		tests.id = ? and results.success = 0
	ORDER BY uploads.time desc
	`, testID)
	if err != nil {
		return nil, err
	}

	failures := []TestFailure{}
	for rows.Next() {
		var failure TestFailure
		var dateTimestamp int64
		err = rows.Scan(&failure.UploadID, &failure.CommitSHA, &failure.URL, &dateTimestamp)
		if err != nil {
			return nil, err
		}
		failure.Date = time.Unix(dateTimestamp, 0)
		failures = append(failures, failure)
	}

	return failures, nil
}
