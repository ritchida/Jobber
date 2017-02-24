package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/ritchida/jobber/generated/jobber/restapi/operations"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/job"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations/jobs"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.JobberAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.JobberAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.TxtConsumer = httpkit.TextConsumer()

	api.TxtProducer = httpkit.TextProducer()

	api.JSONProducer = httpkit.JSONProducer()

	api.JobCreateJobHandler = job.CreateJobHandlerFunc(func(params job.CreateJobParams) middleware.Responder {
		return middleware.NotImplemented("operation job.CreateJob has not yet been implemented")
	})
	api.JobGetJobHandler = job.GetJobHandlerFunc(func(params job.GetJobParams) middleware.Responder {
		return middleware.NotImplemented("operation job.GetJob has not yet been implemented")
	})
	api.JobsGetJobsHandler = jobs.GetJobsHandlerFunc(func(params jobs.GetJobsParams) middleware.Responder {
		return middleware.NotImplemented("operation jobs.GetJobs has not yet been implemented")
	})

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
