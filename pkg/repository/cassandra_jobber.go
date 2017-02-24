package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/ritchida/jobber/pkg/config"
	"github.com/ritchida/jobber/pkg/models"
)

const (
	selectJobsQuery      = "SELECT job_id, created, last_updated, completed, status, tags, type, owner FROM jobs"
	selectJobByIDQuery   = "SELECT job_id, created, last_updated, completed, status, tags, type, owner FROM jobs where job_id = ?"
	insertJobQuery       = "INSERT INTO jobs                (job_id, created, last_updated, status, tags, type, owner) VALUES    (?, ?, ?, ?, ?, ?, ?)"
	insertLatestJobQuery = "INSERT INTO latest_jobs (bucket, job_id, created, last_updated, status, tags, type, owner) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	deleteJobQuery       = "DELETE FROM jobs        where                job_id = ?"
	deleteLatestJobQuery = "DELETE FROM latest_jobs where bucket = ? and job_id = ?"
)

var repository *CassandraJobberRepository
var repoInitOnce sync.Once
var repoInitErr error

// CassandraJobberRepository is an application of the repository pattern for jobs
type CassandraJobberRepository struct {
	JobberRepository
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

// GetCassandraJobberRepository creates an instance of CassandraJobberRepository
func GetCassandraJobberRepository() (*CassandraJobberRepository, error) {
	repoInitOnce.Do(initRepo)
	return repository, repoInitErr
}

func initRepo() {
	repo := CassandraJobberRepository{}

	var configErrors []error
	jobberConfig, configErrors := config.GetJobberConfig()
	if len(configErrors) > 0 {
		for _, err := range configErrors {
			fmt.Printf("Configuration error: %v\n", err)
		}
		return
	}

	// connect to the cluster
	repo.cluster = gocql.NewCluster(jobberConfig.Cassandra.ClusterNodeIPs)
	repo.cluster.Keyspace = "jobber"
	repo.cluster.Consistency = gocql.Quorum

	// connect to the database
	repo.session, repoInitErr = repo.cluster.CreateSession()

	repository = &repo
}

// Close closes the underlying connection to the database
func (r CassandraJobberRepository) Close() {
	r.session.Close()
}

// GetJobs returns all jobs
func (r CassandraJobberRepository) GetJobs() ([]*models.Job, error) {
	var completedAt time.Time
	job := models.Job{}
	jobs := []*models.Job{}
	iter := r.session.Query(selectJobsQuery).Iter()
	for iter.Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &completedAt, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		newJob := models.Job{
			ID:          job.ID,
			CreatedAt:   job.CreatedAt,
			UpdatedAt:   job.UpdatedAt,
			CompletedAt: &completedAt,
			Status:      job.Status,
			Type:        job.Type,
			Tags:        job.Tags,
			Owner:       job.Owner,
		}
		jobs = append(jobs, &newJob)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJob returns the job specified by ID
func (r CassandraJobberRepository) GetJob(ID string) (*models.Job, error) {
	var completedAt time.Time
	job := models.Job{}
	var newJob *models.Job
	iter := r.session.Query(selectJobsQuery).Iter()
	for iter.Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &completedAt, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		newJob = &models.Job{
			ID:          job.ID,
			CreatedAt:   job.CreatedAt,
			UpdatedAt:   job.UpdatedAt,
			CompletedAt: &completedAt,
			Status:      job.Status,
			Type:        job.Type,
			Tags:        job.Tags,
			Owner:       job.Owner,
		}
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return newJob, nil
}

// InsertJob adds the specified job to the job repository
func (r CassandraJobberRepository) InsertJob(job *models.JobSpec) (string, error) {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	timeUUID := gocql.TimeUUID()
	timestamp := timeUUID.Timestamp()
	batch.Query(insertJobQuery, timeUUID, timestamp, timestamp, "queued", job.Tags, job.Type, "danritchie")
	batch.Query(insertLatestJobQuery, 0, timeUUID, timestamp, timestamp, "queued", job.Tags, job.Type, "danritchie")

	return timeUUID.String(), r.session.ExecuteBatch(batch)
}

// DeleteJob removes the specified job from the job repository
func (r CassandraJobberRepository) DeleteJob(ID string) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	batch.Query(deleteJobQuery, ID)
	batch.Query(deleteLatestJobQuery, 0, ID)

	return r.session.ExecuteBatch(batch)
}
