package repository

import (
	"fmt"
	"time"

	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/gocql/gocql"
	"github.com/ritchida/jobber/pkg/models"
)

// CassandraJobberRepository is an application of the repository pattern for jobs
type CassandraJobberRepository struct {
	JobberRepository
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

// NewCassandraJobberRepository creates an instance of CassandraJobberRepository
func NewCassandraJobberRepository() (JobberRepository, error) {
	repo := CassandraJobberRepository{}

	// connect to the cluster
	// TODO: get this from config
	repo.cluster = gocql.NewCluster("35.166.53.200")
	repo.cluster.Keyspace = "jobber"
	repo.cluster.Consistency = gocql.Quorum

	// connect to the database
	var err error
	repo.session, err = repo.cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// Close closes the underlying connection to the database
func (r CassandraJobberRepository) Close() {
	r.session.Close()
}

// GetJobs returns all jobs
func (r CassandraJobberRepository) GetJobs() ([]*models.Job, error) {
	var time1, time2, time3 time.Time
	job := models.Job{}
	jobs := []*models.Job{}
	// iter := r.session.Query(`SELECT * FROM jobs`).Iter()
	iter := r.session.Query(`SELECT job_id, created, last_updated, completed, status, tags, type, owner FROM jobs`).Iter()
	// iter := r.session.Query(`SELECT job_id, status, tags, type, owner FROM jobs`).Iter()
	for _, colInfo := range iter.Columns() {
		fmt.Printf("Name %s, type %s\n", colInfo.Name, colInfo.TypeInfo.Type())
	}
	for iter.Scan(&job.ID, &time1, &time2, &time3, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		// for iter.Scan(&job.ID, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		job.CreatedAt = strfmt.DateTime(time1)
		job.UpdatedAt = strfmt.DateTime(time2)
		if &time3 == nil {
			job.CompletedAt = nil
		} else {
			dt := strfmt.DateTime(time3)
			job.CompletedAt = &dt
		}

		jobs = append(jobs, &job)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJob returns the job specified by ID
func (r CassandraJobberRepository) GetJob(ID string) (*models.Job, error) {
	return nil, nil
}

// InsertJob adds the specified job to the job repository
func (r CassandraJobberRepository) InsertJob(job *models.JobSpec) error {
	return nil
}

// DeleteJob removes the specified job from the job repository
func (r CassandraJobberRepository) DeleteJob(ID string) error {
	return nil
}

func ToDateTime(time *time.Time) {

}
