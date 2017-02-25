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

// GetJobMessages will get a list of compute ids and profiles from rpmgr client
func GetJobMessages(params job.GetJobMessagesParams) middleware.Responder {
	jobRepo, err := repository.GetCassandraJobberRepository()
	if err != nil {
		newErr := fmt.Errorf("Unable to access jobs repository: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewGetJobMessagesDefault(0).WithPayload(&se)
	}

	var list []*models.JobMessage
	list, err = jobRepo.GetJobMessages(params.ID)
	if err != nil {
		newErr := fmt.Errorf("Unable to get job messages: %v", err)
		se := createServiceError(newErr, http.StatusInternalServerError)
		return job.NewGetJobMessagesDefault(0).WithPayload(&se)
	}

	apiJobMessages := []*apimodel.JobMessage{}
	for _, job := range list {
		apiJobMessage := job.ToAPI()
		apiJobMessages = append(apiJobMessages, &apiJobMessage)
	}

	return job.NewGetJobMessagesOK().WithPayload(apiJobMessages)
}
