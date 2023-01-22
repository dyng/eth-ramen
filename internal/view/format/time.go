package format

import "time"

func ToDatetime(sec uint64) string {
	return time.Unix(int64(sec), 0).Format("2006-01-02 15:04:05")
}
