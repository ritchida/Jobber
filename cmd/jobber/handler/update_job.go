package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/pkg/repository"
)

// UpdateJob will Update a new Job with the specified profile
func UpdateJob(params job.UpdateJobParams) middleware.Responder {
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewUpdateJobDefault(0).WithPayload(&se)
	}

	status := string(params.Status)
	switch status {
	case "created":
		return updateJobStatus(jobRepo, params.ID, string(status))
	case "running":
		return updateJobStatus(jobRepo, params.ID, string(status))
	case "unknown":
		return updateJobStatus(jobRepo, params.ID, string(status))
	case "failed":
		return completeJob(jobRepo, params.ID, string(status))
	case "succeeded":
		return completeJob(jobRepo, params.ID, string(status))
	default:
		newErr := fmt.Errorf("Invalid job status specified: '%s'", string(status))
		se := createServiceError(newErr, http.StatusBadRequest)
		return job.NewUpdateJobDefault(0).WithPayload(&se)
	}
}

func updateJobStatus(jobRepo repository.JobberRepository, ID string, status string) middleware.Responder {
	err := jobRepo.UpdateJobStatus(ID, status)
	if err != nil {
		newErr := fmt.Errorf("Unable to update job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewUpdateJobDefault(0).WithPayload(&se)
	}

	return job.NewUpdateJobAccepted()
}

func completeJob(jobRepo repository.JobberRepository, ID string, status string) middleware.Responder {
	err := jobRepo.CompleteJob(ID, status)
	if err != nil {
		newErr := fmt.Errorf("Unable to complete job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewUpdateJobDefault(0).WithPayload(&se)
	}

	return job.NewUpdateJobAccepted()
}
