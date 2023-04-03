package database

import (
	"strconv"

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

func UpdateFlakyTests(db *sqlx.DB, commitID int) error {
	stmt, err := db.Preparex(`
	UPDATE tests SET flaky = true
		WHERE tests.id in (
			SELECT DISTINCT(tests.id) FROM results
			LEFT JOIN tests ON tests.id = results.test_id
			WHERE results.commit_id == ?
			GROUP BY results.test_id
			HAVING COUNT(DISTINCT results.success) > 1
		)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(commitID)
	return err
}

type FlakyTest struct {
	ID        int     `db:"test_id"`
	Name      string  `db:"name"`
	Successes int     `db:"successes"`
	Fails     int     `db:"fails"`
	Score     float32 `db:"score"`
}

func GetFlakyTests(db *sqlx.DB, lastNDays int) ([]FlakyTest, error) {
	rows, err := db.Queryx(`
	SELECT
		tests.id AS test_id,
		tests.name AS name,
		SUM(results.success) AS successes,
		COUNT(*) - SUM(results.success) AS fails,
		1-AVG(results.success) AS score
	FROM results
	JOIN tests ON results.test_id = tests.id
	LEFT JOIN uploads ON uploads.id = results.upload_id
	WHERE
		tests.flaky == true
		AND
			datetime(uploads.time, 'auto') BETWEEN
			datetime('now', ?) AND datetime('now', 'localtime')
	GROUP BY test_id
	HAVING fails > 0
	ORDER BY score DESC
	`, strconv.Itoa(-lastNDays)+` days`)
	if err != nil {
		return nil, err
	}

	tests := []FlakyTest{}
	for rows.Next() {
		var test FlakyTest
		err = rows.StructScan(&test)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

func GetTestName(db *sqlx.DB, testID int) (string, error) {
	name := ""
	err := db.QueryRowx(
		`SELECT tests.name from tests WHERE tests.id = ?
	`, testID).Scan(&name)

	return name, err
}

type TestResult struct {
	UploadID  int     `db:"upload_id"`
	TestID    int     `db:"test_id"`
	Name      string  `db:"test_name"`
	CommitSHA string  `db:"commit_sha"`
	URL       *string `db:"url"`
	Success   bool    `db:"success"`
	Output    *string `db:"test_output"`
	Date      string  `db:"time"`
}

func GetTestResult(db *sqlx.DB, testID, uploadID int) (TestResult, error) {
	output := TestResult{UploadID: uploadID, TestID: testID}
	err := db.QueryRowx(`
	SELECT
		tests.name as test_name,
		commits.commit_sha,
		uploads.url,
		results.success,
		results.output as test_output,
		uploads.time
	from results
	LEFT JOIN commits on commits.id = results.commit_id
	LEFT JOIN tests on tests.id = results.test_id
	LEFT JOIN uploads on uploads.id = results.upload_id
	WHERE
		results.test_id = ? and results.upload_id = ?
	`, testID, uploadID).StructScan(&output)

	return output, err
}
