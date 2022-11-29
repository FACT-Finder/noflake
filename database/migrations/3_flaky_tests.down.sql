DROP INDEX IF EXISTS tests_flaky;
DROP INDEX IF EXISTS results_test_id;
DROP INDEX IF EXISTS results_commit_id;
ALTER TABLE tests DROP flaky;
CREATE INDEX results_idx ON results(test_id, commit_id, success);
