package database

import (
	"github.com/FACT-Finder/noflake/model"
	"github.com/jmoiron/sqlx"
)

func CreateOrGetCommit(db *sqlx.DB, commit model.Commit) (*model.Commit, error) {
	_, err := db.Exec("INSERT OR IGNORE INTO commits (commit_sha) VALUES (?)", commit.CommitSha)
	if err != nil {
		return nil, err
	}

	var id int
	err = db.QueryRowx("SELECT id FROM commits WHERE commit_sha = (?)", commit.CommitSha).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.Commit{ID: &id, CommitSha: commit.CommitSha}, nil
}
