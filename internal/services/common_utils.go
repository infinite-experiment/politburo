package services

import (
	"fmt"
	"time"
)

func GetResponseTime(init time.Time) string {
	timeDiff := time.Since(init).Milliseconds()
	return fmt.Sprintf("%dms", timeDiff)
}
