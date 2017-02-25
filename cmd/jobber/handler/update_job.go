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

	err = jobRepo.UpdateJobStatus(params.ID, string(params.Status))
	if err != nil {
		newErr := fmt.Errorf("Unable to Update job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewUpdateJobDefault(0).WithPayload(&se)
	}

	return job.NewUpdateJobAccepted()
}
