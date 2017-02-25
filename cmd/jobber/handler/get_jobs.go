package handler

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	apimodel "github.com/ritchida/jobber/generated/jobber/models"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/jobs"
	"github.com/ritchida/jobber/pkg/models"
	"github.com/ritchida/jobber/pkg/repository"
)

// GetJobs will get a list of compute ids and profiles from rpmgr client
func GetJobs(params jobs.GetJobsParams) middleware.Responder {
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return jobs.NewGetJobsDefault(0).WithPayload(&se)
	}

	var list []*models.Job
	if params.NumLatest == nil {
		list, err = jobRepo.GetJobs()
		if err != nil {
			newErr := fmt.Errorf("Unable to get jobs: %v", err)
			se := createServiceError(newErr, http.StatusInternalServerError)
			return jobs.NewGetJobsDefault(0).WithPayload(&se)
		}
	} else {
		list, err = jobRepo.GetLatestJobs(int(*params.NumLatest))
		if err != nil {
			newErr := fmt.Errorf("Unable to get jobs: %v", err)
			se := createServiceError(newErr, http.StatusInternalServerError)
			return jobs.NewGetJobsDefault(0).WithPayload(&se)
		}
	}

	apiJobs := []*apimodel.Job{}
	for _, job := range list {
		apiJob := job.ToAPI()
		apiJobs = append(apiJobs, &apiJob)
	}

	return jobs.NewGetJobsOK().WithPayload(apiJobs)
}
