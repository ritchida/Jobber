package models

import (
	"time"

	"github.com/go-swagger/go-swagger/strfmt"
)

// TimeToDateTime converts a Go time.Time to swagger strfmt.DateTime
func TimeToDateTime(time *time.Time) *strfmt.DateTime {
	if time == nil {
		return nil
	}
	dt := strfmt.DateTime(*time)
	return &dt
}

// DateTimeToTime converts a swagger strfmt.DateTime to Go time.Time
func DateTimeToTime(dateTime *strfmt.DateTime) *time.Time {
	if dateTime == nil {
		return nil
	}
	var t time.Time
	t = time.Time(*dateTime)
	return &t
}
