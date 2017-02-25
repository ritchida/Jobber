package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/pkg/repository"
)

// AddJobMessage will create a new Job with the specified profile
func AddJobMessage(params job.AddJobMessageParams) middleware.Responder {
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewCreateJobDefault(0).WithPayload(&se)
	}

	exists, err := jobExists(jobRepo, params.ID)
	if err != nil {
		newErr := fmt.Errorf("Unable to determine if specified job exists: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewCreateJobDefault(0).WithPayload(&se)
	} else if !exists {
		newErr := fmt.Errorf("Unable to locate job: %s", params.ID)
		se := createServiceError(newErr, http.StatusNotFound)
		return job.NewGetJobDefault(0).WithPayload(&se)
	}

	err = jobRepo.AddJobMessage(params.ID, params.Message)
	if err != nil {
		newErr := fmt.Errorf("Unable to add job message: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewCreateJobDefault(0).WithPayload(&se)
	}

	return job.NewCreateJobAccepted()
}

func jobExists(jobRepo repository.JobberRepository, ID string) (bool, error) {
	j, err := jobRepo.GetJob(ID)
	if err != nil {
		return false, err
	}

	if j == nil {
		return false, nil
	}

	return true, nil

}
