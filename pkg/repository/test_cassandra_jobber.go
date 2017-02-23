package repository

import (
	"github.com/gocql/gocql"
	"github.com/ritchida/jobber/pkg/models"
)

// TestJobberRepository is a test interface following the repository pattern for storing jobs
type TestJobberRepository interface {
	JobberRepository
	TruncateTables() []error
}

// TestCassandraJobberRepository is a test application of the repository pattern for jobs
type TestCassandraJobberRepository struct {
	TestJobberRepository
	repo CassandraJobberRepository
}

// NewTestCassandraJobberRepository creates an instance of TestCassandraJobberRepository
func NewTestCassandraJobberRepository() (TestJobberRepository, error) {
	repo, err := NewCassandraJobberRepository()
	if err != nil {
		return nil, err
	}
	testRepo := TestCassandraJobberRepository{
		repo: *repo,
	}

	return testRepo, nil
}

// TruncateTables returns all jobs
func (tr TestCassandraJobberRepository) TruncateTables() []error {
	var err error
	errors := []error{}
	var query *gocql.Query
	query = tr.repo.session.Query(`TRUNCATE jobs`)
	err = query.Exec()
	if err != nil {
		errors = append(errors, err)
	}
	query = tr.repo.session.Query(`TRUNCATE latest_jobs`)
	err = query.Exec()
	if err != nil {
		errors = append(errors, err)
	}
	return errors
}

// Close closes the underlying connection to the database
func (tr TestCassandraJobberRepository) Close() {
	tr.repo.Close()
}

// GetJobs returns all jobs
func (tr TestCassandraJobberRepository) GetJobs() ([]*models.Job, error) {
	return tr.repo.GetJobs()
}

// GetJob returns the job specified by ID
func (tr TestCassandraJobberRepository) GetJob(ID string) (*models.Job, error) {
	return tr.repo.GetJob(ID)
}

// InsertJob adds the specified job to the job repository
func (tr TestCassandraJobberRepository) InsertJob(job *models.JobSpec) (string, error) {
	return tr.repo.InsertJob(job)
}

// DeleteJob removes the specified job from the job repository
func (tr TestCassandraJobberRepository) DeleteJob(ID string) error {
	return tr.repo.DeleteJob(ID)
}
