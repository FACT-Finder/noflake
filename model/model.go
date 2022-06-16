package model

import "time"

type Test struct {
	Name string
	ID   *int
}

type Commit struct {
	CommitSha string
	ID        *int
}

type Upload struct {
	CommitID int
	Time     time.Time
	URL      *string
	ID       *int
}

type Result struct {
	TestID   int
	UploadID int
	CommitID int
	Success  bool
	Output   *string
}

type TestResult struct {
	Name    string
	Success bool
	Output  *string
}
