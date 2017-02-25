package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/ritchida/jobber/pkg/models"
	"github.com/stretchr/testify/assert"
)

// var jobRepo TestJobberRepository
var jobRepo *CassandraJobberRepository

func TestMain(m *testing.M) {
	var err error
	jobRepo, err = GetCassandraJobberRepository()
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	returnCode := m.Run()

	jobRepo.Close()

	os.Exit(returnCode)
}

func TestJobLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	assert := assert.New(t)

	var err error
	var jobs []*models.Job
	var jobID string

	jobSpec := models.JobSpec{
		Type: "integration test",
	}

	// insert 3 jobs, and make sure we can query each by its ID
	insertedJobs := []*models.Job{}
	for count := 0; count < 3; count++ {
		jobID, err = jobRepo.InsertJob(&jobSpec)
		assert.NoError(err)
		assert.NotNil(jobID)

		jobs, err = jobRepo.GetJobs()
		assert.NoError(err)
		assert.NotEmpty(jobs)

		var job *models.Job
		for _, j := range jobs {
			if j.ID == jobID {
				job = j
				break
			}
		}

		jobByID, err := jobRepo.GetJob(jobID)
		assert.NoError(err)
		assert.NotNil(jobByID)
		assertEqualJobs(t, job, jobByID)

		insertedJobs = append(insertedJobs, jobByID)
	}

	// update the status on all 3 jobs, and verify status and updated time
	updatedJobs := []*models.Job{}
	for _, j := range insertedJobs {
		jobRepo.UpdateJobStatus(j.ID, "running")
		jobByID, err := jobRepo.GetJob(j.ID)
		assert.NoError(err)
		assert.Equal("running", jobByID.Status)
		assert.NotEqual(j.UpdatedAt, jobByID.UpdatedAt)
		updatedJobs = append(updatedJobs, jobByID)
	}

	// complete all 3 jobs and verify status and completion time
	completedJobs := []*models.Job{}
	for _, j := range updatedJobs {
		jobRepo.CompleteJob(j.ID, "completed")
		jobByID, err := jobRepo.GetJob(j.ID)
		assert.NoError(err)
		assert.Equal("completed", jobByID.Status)
		assert.NotEqual(j.UpdatedAt, jobByID.UpdatedAt)
		assert.NotEqual(j.CompletedAt, jobByID.CompletedAt)
		completedJobs = append(completedJobs, jobByID)
	}

	// make sure we can get the latest jobs correctly
	latestJobs, err := jobRepo.GetLatestJobs(2)
	assert.NoError(err)
	assert.NotEmpty(latestJobs)
	assert.Equal(2, len(latestJobs))
	assertEqualJobs(t, completedJobs[2], latestJobs[0])
	assertEqualJobs(t, completedJobs[1], latestJobs[1])

	// delete the jobs we ceated, and make sure we can no longer query them by ID
	for _, job := range insertedJobs {
		err = jobRepo.DeleteJob(job.ID)
		assert.NoError(err)
		jobByID, err := jobRepo.GetJob(job.ID)
		assert.NoError(err)
		assert.Nil(jobByID)
	}
}

func assertEqualJobs(t *testing.T, expectedJob *models.Job, actualJob *models.Job) {
	assert := assert.New(t)
	assert.Equal(expectedJob.ID, actualJob.ID)
	assert.Equal(expectedJob.CreatedAt, actualJob.CreatedAt)
	assert.Equal(expectedJob.UpdatedAt, actualJob.UpdatedAt)
	assert.Equal(expectedJob.CompletedAt, actualJob.CompletedAt)
	assert.Equal(expectedJob.Status, actualJob.Status)
	assert.Equal(expectedJob.Type, actualJob.Type)
	assert.Equal(expectedJob.Tags, actualJob.Tags)
	assert.Equal(expectedJob.Owner, actualJob.Owner)
}

func printJobs(jobs []models.Job) {
	for _, job := range jobs {
		fmt.Printf("Job id: %#v\n", job.ID)
		fmt.Printf("Job created at: %#v\n", job.CreatedAt)
		fmt.Printf("Job updated at: %#v\n", job.UpdatedAt)
		if job.CompletedAt != nil {
			fmt.Println("completedAt NOT NIL!!!")
		} else {
			fmt.Println("completedAt NIL!!!")
		}
		fmt.Printf("Job completed at: %#v\n", job.CompletedAt)
		fmt.Printf("Job status: %#v\n", job.Status)
		fmt.Printf("Job type: %#v\n", job.Type)
		fmt.Printf("Job tags: %#v\n", job.Tags)
		fmt.Printf("Job owner: %#v\n", job.Owner)
	}
}
