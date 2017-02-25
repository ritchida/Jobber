package models

import (
	"time"

	"github.com/ritchida/jobber/generated/jobber/models"
)

// JobMessage is an internal representation of the API JobMessage type
type JobMessage struct {
	CreatedAt time.Time
	Message   string
}

// ToAPI converts an internal JobMessage object to an API JobMessage object
func (t *JobMessage) ToAPI() models.JobMessage {
	return models.JobMessage{
		CreatedAt: *TimeToDateTime(&t.CreatedAt),
		Message:   t.Message,
	}
}
