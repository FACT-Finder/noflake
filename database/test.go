package database

import (
	"time"

	"github.com/FACT-Finder/noflake/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func InsertTests(db *sqlx.DB, tests []model.TestResult, upload model.Upload) error {
	stmt, err := db.Preparex("INSERT OR IGNORE INTO tests (name) VALUES (?)")
	if err != nil {
		return err
	}

	for _, test := range tests {
		_, err = stmt.Exec(test.Name)
		if err != nil {
			return err
		}
	}

	stmt, err = db.Preparex("SELECT id FROM tests WHERE name = ?")
	if err != nil {
		return err
	}

	duplicates := map[string]bool{}
	results := map[int]model.Result{}
	var testID int
	for _, test := range tests {
		err = stmt.QueryRowx(test.Name).Scan(&testID)
		if err != nil {
			return err
		}

		if result, exists := results[testID]; exists {
			duplicates[test.Name] = true
			if !test.Success {
				result.Success = false
			}
		} else {
			results[testID] = model.Result{
				TestID:   testID,
				UploadID: *upload.ID,
				CommitID: upload.CommitID,
				Success:  test.Success,
				Output:   test.Output,
			}
		}
	}

	stmt, err = db.Preparex(
		"INSERT INTO results (test_id, upload_id, commit_id, success, output) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	for _, result := range results {
		_, err = stmt.Exec(
			result.TestID, result.UploadID, result.CommitID, result.Success, result.Output)
		if err != nil {
			log.Err(err).
				Int("test", result.TestID).
				Int("upload", result.UploadID).
				Int("commit", result.CommitID).
				Msg("couldn't insert into results")
			return err
		}
	}

	if len(duplicates) > 0 {
		names := make([]string, 0, len(duplicates))
		for k := range duplicates {
			names = append(names, k)
		}
		log.Warn().Strs("tests", names).Msg("upload contains duplicate test names")
	}

	return nil
}

func GetFlakyTests(db *sqlx.DB) ([]FlakyTest, error) {
	rows, err := db.Queryx(`
	WITH flaky_tests (id, name) AS (
		SELECT DISTINCT(tests.id), tests.name FROM results
		LEFT JOIN tests ON tests.id = results.test_id
		GROUP BY results.test_id, results.commit_id
		HAVING COUNT(DISTINCT results.success) > 1
	)
	SELECT
		flaky_tests.name AS name,
		count(results.success) as total_fails,
		MAX(uploads.time) as last_fail
	FROM results
	JOIN flaky_tests ON results.test_id = flaky_tests.id
	LEFT JOIN uploads ON uploads.id = results.upload_id
	WHERE results.success == false
	GROUP BY test_id
	`)
	if err != nil {
		return nil, err
	}

	tests := []FlakyTest{}
	for rows.Next() {
		var test FlakyTest
		var lastFailTimestamp int64
		err = rows.Scan(&test.Name, &test.TotalFails, &lastFailTimestamp)
		if err != nil {
			return nil, err
		}
		test.LastFail = time.Unix(lastFailTimestamp, 0)
		tests = append(tests, test)
	}

	return tests, nil
}

type FlakyTest struct {
	Name       string    `db:"name"`
	TotalFails int       `db:"total_fails"`
	LastFail   time.Time `db:"last_fail"`
}
