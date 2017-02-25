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
	selectLatestJobsQuery           = "SELECT job_id, created, last_updated, completed, status, tags, type, owner FROM latest_jobs where bucket = ? LIMIT ?"
	selectJobsQuery                 = "SELECT job_id, created, last_updated, completed, status, tags, type, owner FROM jobs"
	selectJobByIDQuery              = "SELECT job_id, created, last_updated, completed, status, tags, type, owner FROM jobs where job_id = ?"
	insertJobQuery                  = "INSERT INTO jobs                (job_id, created, last_updated, status, tags, type, owner) VALUES    (?, ?, ?, ?, ?, ?, ?)"
	insertLatestJobQuery            = "INSERT INTO latest_jobs (bucket, job_id, created, last_updated, status, tags, type, owner) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	updateJobStatusQuery            = "UPDATE        jobs SET status = ?, last_updated = ? where                job_id = ?"
	updateLatestJobStatusQuery      = "UPDATE latest_jobs SET status = ?, last_updated = ? where bucket = ? and job_id = ?"
	completeJobQuery                = "UPDATE        jobs SET status = ?, last_updated = ?, completed = ? where                job_id = ?"
	completeLatestJobQuery          = "UPDATE latest_jobs SET status = ?, last_updated = ?, completed = ? where bucket = ? and job_id = ?"
	addJobMessageQuery              = "INSERT into job_messages_by_job_id (job_id, message_created, message) values (?, ?, ?)"
	updateJobLastUpdatedQuery       = "UPDATE        jobs SET last_updated = ? where                job_id = ?"
	updateLatestJobLastUpdatedQuery = "UPDATE latest_jobs SET last_updated = ? where bucket = ? and job_id = ?"
	getJobMessagesQuery             = "SELECT message_created, message FROM job_messages_by_job_id where job_id = ?"
	deleteJobQuery                  = "DELETE FROM jobs        where                job_id = ?"
	deleteLatestJobQuery            = "DELETE FROM latest_jobs where bucket = ? and job_id = ?"
	deleteJobMessagesQuery          = "DELETE FROM job_messages_by_job_id where job_id = ?"
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

// GetLatestJobs returns the "numJobs" most recently created jobs
func (r CassandraJobberRepository) GetLatestJobs(numJobs int) ([]*models.Job, error) {
	var completedAt time.Time
	job := models.Job{}
	jobs := []*models.Job{}
	iter := r.session.Query(selectLatestJobsQuery, 0, numJobs).Iter()
	for iter.Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &completedAt, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		completedAtValue := completedAt
		newJob := models.Job{
			ID:          job.ID,
			CreatedAt:   job.CreatedAt,
			UpdatedAt:   job.UpdatedAt,
			CompletedAt: &completedAtValue,
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

// GetJobs returns all jobs
func (r CassandraJobberRepository) GetJobs() ([]*models.Job, error) {
	var completedAt time.Time
	job := models.Job{}
	jobs := []*models.Job{}
	iter := r.session.Query(selectJobsQuery).Iter()
	for iter.Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &completedAt, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		completedAtValue := completedAt
		newJob := models.Job{
			ID:          job.ID,
			CreatedAt:   job.CreatedAt,
			UpdatedAt:   job.UpdatedAt,
			CompletedAt: &completedAtValue,
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
	iter := r.session.Query(selectJobByIDQuery, ID).Iter()
	for iter.Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &completedAt, &job.Status, &job.Tags, &job.Type, &job.Owner) {
		completedAtValue := completedAt
		newJob = &models.Job{
			ID:          job.ID,
			CreatedAt:   job.CreatedAt,
			UpdatedAt:   job.UpdatedAt,
			CompletedAt: &completedAtValue,
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
	time := timeUUID.Time()
	batch.Query(insertJobQuery, timeUUID, time, time, "queued", job.Tags, job.Type, "danritchie")
	batch.Query(insertLatestJobQuery, 0, timeUUID, time, time, "queued", job.Tags, job.Type, "danritchie")

	return timeUUID.String(), r.session.ExecuteBatch(batch)
}

// UpdateJobStatus updates the specified job's status
func (r CassandraJobberRepository) UpdateJobStatus(ID string, status string) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	timeUUID := gocql.TimeUUID()
	time := timeUUID.Time()
	batch.Query(updateJobStatusQuery, status, time, ID)
	batch.Query(updateLatestJobStatusQuery, status, time, 0, ID)

	return r.session.ExecuteBatch(batch)
}

// AddJobMessage adds a timestamped message to the specified job
func (r CassandraJobberRepository) AddJobMessage(ID string, message string) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	timeUUID := gocql.TimeUUID()
	time := timeUUID.Time()
	batch.Query(addJobMessageQuery, ID, time, message)
	batch.Query(updateJobLastUpdatedQuery, time, ID)
	batch.Query(updateLatestJobLastUpdatedQuery, time, 0, ID)

	return r.session.ExecuteBatch(batch)
}

// GetJobMessages adds a timestamped message to the specified job
func (r CassandraJobberRepository) GetJobMessages(ID string) ([]*models.JobMessage, error) {
	jobMsg := models.JobMessage{}
	jobMsgs := []*models.JobMessage{}
	iter := r.session.Query(getJobMessagesQuery, ID).Iter()
	for iter.Scan(&jobMsg.CreatedAt, &jobMsg.Message) {
		newJobMsg := &models.JobMessage{
			CreatedAt: jobMsg.CreatedAt,
			Message:   jobMsg.Message,
		}
		jobMsgs = append(jobMsgs, newJobMsg)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return jobMsgs, nil
}

// CompleteJob updates the specified job's status
func (r CassandraJobberRepository) CompleteJob(ID string, status string) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	timeUUID := gocql.TimeUUID()
	time := timeUUID.Time()
	batch.Query(completeJobQuery, status, time, time, ID)
	batch.Query(completeLatestJobQuery, status, time, time, 0, ID)

	return r.session.ExecuteBatch(batch)
}

// DeleteJob removes the specified job from the job repository
func (r CassandraJobberRepository) DeleteJob(ID string) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	batch.Query(deleteJobQuery, ID)
	batch.Query(deleteLatestJobQuery, 0, ID)

	return r.session.ExecuteBatch(batch)
}
