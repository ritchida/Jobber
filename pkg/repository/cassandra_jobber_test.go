package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJobberRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	assert := assert.New(t)

	jobRepo, err := NewCassandraJobberRepository()
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	defer jobRepo.Close()

	jobs, err := jobRepo.GetJobs()
	assert.NoError(err)
	assert.Empty(jobs)

	// jobSpec := models.JobSpec{
	// 	Type: "integration test",
	// }

	// err = jobRepo.InsertJob(&jobSpec)
	// assert.NoError(err)

	// jobs, err = jobRepo.GetJobs()
	// assert.NoError(err)
	// assert.NotEmpty(jobs)

	// job := jobs[0]

	// jobByID, err := jobRepo.GetJob(job.ID)
	// assert.NoError(err)
	// assert.Equal(job.ID, jobByID.ID)
	// assert.Equal(job.CreatedAt, jobByID.CreatedAt)
	// assert.Equal(job.UpdatedAt, jobByID.UpdatedAt)
	// assert.Equal(job.UpdatedAt, jobByID.UpdatedAt)
	// assert.Equal(job.CompletedAt, jobByID.CompletedAt)
	// assert.Equal(job.Status, jobByID.Status)
	// assert.Equal(job.Type, jobByID.Type)
	// assert.Equal(job.Owner, jobByID.Owner)

	// err = jobRepo.DeleteTask(task.ID)
	// assert.NoError(err)
}
