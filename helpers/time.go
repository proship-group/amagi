package helpers

import (
	"fmt"
	"time"
)

// TimeToStr convert time to date
func TimeToStr(t time.Time) string {
	// return time.Date(s.Year(), s.Month(), s.Day(), hour, min, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	formatted := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	return formatted.Format(layout)
}

// TimeByZone load time from specific zone
func TimeByZone(zone string) time.Time {
	t := time.Now()
	locationTime, err := time.LoadLocation(zone)
	if err != nil {
		return t
	}

	fmt.Println("time tokyo", t.In(locationTime))
	return t.In(locationTime)
}

// GetTimeFromTokyo get time from tokyo
func GetTimeFromTokyo() time.Time {
	return TimeByZone("Asia/Tokyo")
}
