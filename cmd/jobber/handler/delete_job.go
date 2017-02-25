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
	fmt.Printf("Delete job handler, ID: %s\n", params.ID)
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewGetJobDefault(0).WithPayload(&se)
	}

	fmt.Printf("Deleting job %s from repository\n", params.ID)
	err = jobRepo.DeleteJob(params.ID)
	if err != nil {
		newErr := fmt.Errorf("Unable to delete job: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewDeleteJobDefault(0).WithPayload(&se)
	}

	fmt.Printf("Deleted job %s from repository, returning\n", params.ID)
	return job.NewGetJobOK()
}
