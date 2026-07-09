package github

import "time"

type Commit struct {
	Repository string
	AuthorDate string
}

type Activity struct {
	CommitsToday int
	CommitsWeek  int
	LastCommitAt *time.Time
}
