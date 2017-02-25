package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	apimodel "github.com/ritchida/jobber/generated/jobber/models"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/pkg/models"
	"github.com/ritchida/jobber/pkg/repository"
)

// CreateJob will create a new Job with the specified profile
func CreateJob(params job.CreateJobParams) middleware.Responder {
	if params.JobSpec.Type == "" {
		params.JobSpec.Type = "default"
	}

	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewCreateJobDefault(0).WithPayload(&se)
	}

	internalJobSpec := models.JobSpec{}
	internalJobSpec.FromAPI(*params.JobSpec)
	id, err := jobRepo.InsertJob(&internalJobSpec)
	if err != nil {
		newErr := fmt.Errorf("Unable to create job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewCreateJobDefault(0).WithPayload(&se)
	}

	return job.NewCreateJobAccepted().WithPayload(apimodel.ID(id))
}
