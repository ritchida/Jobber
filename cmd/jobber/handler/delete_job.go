package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/pkg/repository"
)

// DeleteJob deletes a Job with the specified ID.
func DeleteJob(params job.DeleteJobParams) middleware.Responder {
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewGetJobDefault(0).WithPayload(&se)
	}

	err = jobRepo.DeleteJob(params.ID)
	if err != nil {
		newErr := fmt.Errorf("Unable to delete job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewDeleteJobDefault(0).WithPayload(&se)
	}

	return job.NewGetJobOK()
}
