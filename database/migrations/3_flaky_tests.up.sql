DROP INDEX IF EXISTS results_idx;

ALTER TABLE tests ADD flaky INTEGER NOT NULL DEFAULT 0;

UPDATE tests SET flaky = true
    WHERE tests.id in (
        SELECT DISTINCT(tests.id) FROM results
        LEFT JOIN tests ON tests.id = results.test_id
        GROUP BY results.test_id, results.commit_id
        HAVING COUNT(DISTINCT results.success) > 1
    );

CREATE INDEX results_test_id ON results(test_id);
CREATE INDEX results_commit_id ON results(commit_id);
CREATE INDEX tests_flaky ON tests(flaky);
