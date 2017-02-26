package restapi

import (
	"fmt"
	"os"
	"testing"

	httptransport "github.com/go-swagger/go-swagger/httpkit/client"
	"github.com/go-swagger/go-swagger/strfmt"
	jobclient "github.com/ritchida/jobber/generated/jobber-client/client/job"
	jobsclient "github.com/ritchida/jobber/generated/jobber-client/client/jobs"
	"github.com/ritchida/jobber/generated/jobber-client/models"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	returnCode := m.Run()
	os.Exit(returnCode)
}

func TestJobLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	assert := assert.New(t)

	jobClient := jobclient.New(httptransport.New("localhost", "/", []string{"http"}), strfmt.Default)
	jobsClient := jobsclient.New(httptransport.New("localhost", "/", []string{"http"}), strfmt.Default)

	jobSpec := models.JobSpec{
		Type: "integration test",
	}
	params := jobclient.NewCreateJobParams().WithJobSpec(&jobSpec)
	// insert 3 jobs, and make sure we can query each by its ID
	insertedJobs := []*models.Job{}
	for count := 0; count < 3; count++ {
		accepted, err := jobClient.CreateJob(params)
		assert.NoError(err)
		assert.NotNil(accepted)
		assert.NotEmpty(string(accepted.Payload))

		jobsOK, err := jobsClient.GetJobs(jobsclient.NewGetJobsParams())
		assert.NoError(err)
		assert.NotNil(jobsOK)

		var job *models.Job
		for _, j := range jobsOK.Payload {
			if j.ID == accepted.Payload {
				job = j
				break
			}
		}

		jobOK, err := jobClient.GetJob(jobclient.NewGetJobParams().WithID(string(job.ID)))
		assert.NoError(err)
		assert.NotNil(jobOK.Payload)
		assertEqualJobs(t, job, jobOK.Payload)

		insertedJobs = append(insertedJobs, jobOK.Payload)
	}

	// update the status on all 3 jobs, and verify status and updated time
	updatedJobs := []*models.Job{}
	for _, j := range insertedJobs {
		params := jobclient.UpdateJobParams{
			ID:     string(j.ID),
			Status: "running",
		}
		jobClient.UpdateJob(&params)
		jobOK, err := jobClient.GetJob(jobclient.NewGetJobParams().WithID(string(j.ID)))
		assert.NoError(err)
		assert.NotNil(jobOK.Payload)
		assert.Equal("running", jobOK.Payload.Status)
		assert.NotEqual(j.UpdatedAt, jobOK.Payload.UpdatedAt)
		updatedJobs = append(updatedJobs, jobOK.Payload)
	}

	// Add messages to each job and verify the messages can be queried
	for idx, j := range updatedJobs {
		for x := 0; x <= idx; x++ {
			params := jobclient.AddJobMessageParams{
				ID:      string(j.ID),
				Message: fmt.Sprintf("message-%d", x),
			}
			jobClient.AddJobMessage(&params)
		}
		jobOK, err := jobClient.GetJob(jobclient.NewGetJobParams().WithID(string(j.ID)))
		assert.NoError(err)
		assert.NotNil(jobOK.Payload)
		assert.NotEqual(j.UpdatedAt, jobOK.Payload.UpdatedAt)
		jobMsgsOK, err := jobClient.GetJobMessages(jobclient.NewGetJobMessagesParams().WithID(string(j.ID)))
		assert.NoError(err)
		for x := 0; x <= idx; x++ {
			assert.Equal(jobMsgsOK.Payload[x].Message, fmt.Sprintf("message-%d", x))
		}
	}

	// complete all 3 jobs and verify status and completion time
	completedJobs := []*models.Job{}
	for _, j := range updatedJobs {
		params := jobclient.UpdateJobParams{
			ID:     string(j.ID),
			Status: "completed",
		}
		accepted, err := jobClient.UpdateJob(&params)
		assert.NoError(err)
		assert.NotNil(accepted)
		jobOK, err := jobClient.GetJob(jobclient.NewGetJobParams().WithID(string(j.ID)))
		assert.NoError(err)
		assert.NotNil(jobOK.Payload)
		assert.Equal("completed", jobOK.Payload.Status)
		assert.NotEqual(j.UpdatedAt, jobOK.Payload.UpdatedAt)
		assert.NotEqual(j.CompletedAt, jobOK.Payload.CompletedAt)
		assert.Equal(jobOK.Payload.UpdatedAt, jobOK.Payload.CompletedAt)
		completedJobs = append(completedJobs, jobOK.Payload)
	}

	// make sure we can get the latest jobs correctly
	numLatest := int64(2)
	jobsOK, err := jobsClient.GetJobs(jobsclient.NewGetJobsParams().WithNumLatest(&numLatest))
	assert.NoError(err)
	assert.NotEmpty(jobsOK.Payload)
	assert.Equal(2, len(jobsOK.Payload))
	assertEqualJobs(t, completedJobs[2], jobsOK.Payload[0])
	assertEqualJobs(t, completedJobs[1], jobsOK.Payload[1])

	// delete the jobs we ceated, and make sure we can no longer query them by ID
	for _, j := range insertedJobs {
		deleteJobOK, err := jobClient.DeleteJob(jobclient.NewDeleteJobParams().WithID(string(j.ID)))
		assert.NoError(err)
		assert.NotNil(deleteJobOK)
		jobOK, err := jobClient.GetJob(jobclient.NewGetJobParams().WithID(string(j.ID)))
		assert.NoError(err)
		assert.NotNil(jobOK.Payload)
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
