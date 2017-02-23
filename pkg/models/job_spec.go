package models

import "github.com/ritchida/jobber/generated/jobber/models"

// JobSpec is an internal representation of the API JobSpec type
type JobSpec struct {
	Type string
	Tags []string
}

// FromAPI creates a JobSpec from an API JobSpec object
func (t *JobSpec) FromAPI(apiJob models.JobSpec) {
	t.Type = apiJob.Type
	tags := []string{}
	for _, tag := range apiJob.Tags {
		tags = append(tags, string(tag))
	}
	t.Tags = tags
}
