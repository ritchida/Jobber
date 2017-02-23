package models

import (
	"time"

	strfmt "github.com/go-swagger/go-swagger/strfmt"
	"github.com/ritchida/jobber/generated/jobber/models"
)

// Job is an internal representation of the API Job type
type Job struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	Status      string
	Type        string
	Tags        []string
	Owner       string
}

// FromAPI creates a Job from an API Job object
func (t *Job) FromAPI(apiJob models.Job) {
	t.ID = string(apiJob.ID)
	t.CreatedAt = *DateTimeToTime(&apiJob.CreatedAt)
	t.UpdatedAt = *DateTimeToTime(&apiJob.UpdatedAt)
	t.CompletedAt = DateTimeToTime(apiJob.CompletedAt)
	t.Status = apiJob.Status
	t.Type = apiJob.Type
	tags := []string{}
	for _, tag := range apiJob.Tags {
		tags = append(tags, string(tag))
	}
	t.Tags = tags
	t.Owner = apiJob.Owner
}

// ToAPI converts an internal Job object to an API Job object
func (t *Job) ToAPI() models.Job {
	tags := []models.Tag{}
	for _, tag := range t.Tags {
		tags = append(tags, models.Tag(tag))
	}
	return models.Job{
		ID:          models.ID(t.ID),
		CreatedAt:   *TimeToDateTime(&t.CreatedAt),
		UpdatedAt:   *TimeToDateTime(&t.UpdatedAt),
		CompletedAt: TimeToDateTime(t.CompletedAt),
		Status:      t.Status,
		Type:        t.Type,
		Tags:        tags,
		Owner:       t.Owner,
	}
}

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
