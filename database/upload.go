package database

import (
	"github.com/FACT-Finder/noflake/model"
	"github.com/jmoiron/sqlx"
)

func CreateUpload(db *sqlx.DB, upload model.Upload) (*model.Upload, error) {
	var id int
	err := db.QueryRowx(
		"INSERT INTO uploads (commit_id, time, url) VALUES (?, ?, ?) RETURNING (id)",
		upload.CommitID, upload.Time.Unix(), upload.URL).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.Upload{
		ID:       &id,
		CommitID: upload.CommitID,
		Time:     upload.Time,
		URL:      upload.URL,
	}, nil
}
