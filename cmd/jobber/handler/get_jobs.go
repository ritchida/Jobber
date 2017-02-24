package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	apimodel "github.com/ritchida/jobber/generated/jobber/models"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/jobs"
	"github.com/ritchida/jobber/pkg/config"
	"github.com/ritchida/jobber/pkg/repository"
)

// GetJobs will get a list of compute ids and profiles from rpmgr client
func GetJobs(params jobs.GetJobsParams) middleware.Responder {
	// TODO: re-use config
	config, errors := config.ReadConfig()
	for _, configErr := range errors {
		se := apimodel.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%v", configErr),
		}
		return jobs.NewGetJobsDefault(0).WithPayload(&se)
	}

	// TODO: re-use repository
	jobRepo, err := repository.NewCassandraJobberRepository(config.Cassandra)
	if err != nil {
		newErr := fmt.Errorf("Unable to connect to database: %v", err)
		se := apimodel.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%v", newErr),
		}
		return jobs.NewGetJobsDefault(0).WithPayload(&se)
	}

	list, err := jobRepo.GetJobs()
	if err != nil {
		newErr := fmt.Errorf("Unable to connect to database: %v", err)
		se := apimodel.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%v", newErr),
		}
		return jobs.NewGetJobsDefault(0).WithPayload(&se)
	}

	apiJobs := []*apimodel.Job{}
	for _, job := range list {
		apiJob := job.ToAPI()
		apiJobs = append(apiJobs, &apiJob)
	}

	return jobs.NewGetJobsOK().WithPayload(apiJobs)
}
