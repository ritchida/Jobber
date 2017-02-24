package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/ritchida/jobber/pkg/models"
	"github.com/stretchr/testify/assert"
)

var jobRepo TestJobberRepository

func TestMain(m *testing.M) {
	var err error
	jobRepo, err = GetTestCassandraJobberRepository()
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
	assert.Equal(job.ID, jobByID.ID)
	assert.Equal(job.CreatedAt, jobByID.CreatedAt)
	assert.Equal(job.UpdatedAt, jobByID.UpdatedAt)
	assert.Equal(job.UpdatedAt, jobByID.UpdatedAt)
	assert.Equal(job.CompletedAt, jobByID.CompletedAt)
	assert.Equal(job.Status, jobByID.Status)
	assert.Equal(job.Type, jobByID.Type)
	assert.Equal(job.Owner, jobByID.Owner)

	err = jobRepo.DeleteJob(job.ID)
	assert.NoError(err)
	jobByID, err = jobRepo.GetJob(job.ID)
	assert.NoError(err)
	assert.Nil(jobByID)
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
