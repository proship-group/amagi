package helpers

import "time"

// TimeToStr convert time to date
func TimeToStr(t time.Time) string {
	// return time.Date(s.Year(), s.Month(), s.Day(), hour, min, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	formatted := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	return formatted.Format(layout)
}
