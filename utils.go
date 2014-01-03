package scheduler

import "time"

func ToUnixTimestamp(t time.Time) int64 {
	return t.UTC().Truncate(time.Second).Unix()
}
