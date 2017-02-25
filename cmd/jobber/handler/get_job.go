package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/pkg/repository"
)

// GetJob returns a response with a Job payload where the Job ID is "params.ID".
func GetJob(params job.GetJobParams) middleware.Responder {
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewGetJobDefault(0).WithPayload(&se)
	}

	j, err := jobRepo.GetJob(params.ID)
	if err != nil {
		newErr := fmt.Errorf("Unable to get job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewGetJobDefault(0).WithPayload(&se)
	}

	apiJob := j.ToAPI()
	return job.NewGetJobOK().WithPayload(&apiJob)
}
