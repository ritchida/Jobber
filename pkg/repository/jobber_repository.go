package repository

import "github.com/ritchida/jobber/pkg/models"

// JobberRepository is an application of the repository pattern for storing jobs
type JobberRepository interface {
	GetLatestJobs(numJobs int) ([]*models.Job, error)
	GetJobs() ([]*models.Job, error)
	GetJob(ID string) (*models.Job, error)
	InsertJob(job *models.JobSpec) (string, error)
	DeleteJob(ID string) error
	Close()
}
