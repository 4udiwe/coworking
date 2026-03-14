package analytics_service

import "errors"

var (
	ErrCannotInsertEvents = errors.New("cannot insert events")
	ErrCannotFetchInfo    = errors.New("cannot fetch info")
)
