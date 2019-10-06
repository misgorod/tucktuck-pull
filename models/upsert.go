package models

import "time"

type UpsertResult struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    interface{}
	LastTime      time.Time
}
