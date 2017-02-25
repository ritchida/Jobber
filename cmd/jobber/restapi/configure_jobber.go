package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"

	"github.com/ritchida/jobber/cmd/jobber/handler"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/jobs"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.JobberAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

// ConfigureAPI Configures the handlers to serve the jobber API
func configureAPI(api *operations.JobberAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.TxtConsumer = httpkit.TextConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.TxtProducer = httpkit.TextProducer()

	// Add single job handlers here
	api.JobCreateJobHandler = job.CreateJobHandlerFunc(handler.CreateJob)
	api.JobUpdateJobHandler = job.UpdateJobHandlerFunc(handler.UpdateJob)
	api.JobDeleteJobHandler = job.DeleteJobHandlerFunc(handler.DeleteJob)
	api.JobGetJobHandler = job.GetJobHandlerFunc(handler.GetJob)
	api.JobAddJobMessageHandler = job.AddJobMessageHandlerFunc(handler.AddJobMessage)
	api.JobGetJobMessagesHandler = job.GetJobMessagesHandlerFunc(handler.GetJobMessages)

	// add multiple job handlers here
	api.JobsGetJobsHandler = jobs.GetJobsHandlerFunc(handler.GetJobs)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
