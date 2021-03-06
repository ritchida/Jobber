package job

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
)

// UpdateJobHandlerFunc turns a function with the right signature into a update job handler
type UpdateJobHandlerFunc func(UpdateJobParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateJobHandlerFunc) Handle(params UpdateJobParams) middleware.Responder {
	return fn(params)
}

// UpdateJobHandler interface for that can handle valid update job params
type UpdateJobHandler interface {
	Handle(UpdateJobParams) middleware.Responder
}

// NewUpdateJob creates a new http.Handler for the update job operation
func NewUpdateJob(ctx *middleware.Context, handler UpdateJobHandler) *UpdateJob {
	return &UpdateJob{Context: ctx, Handler: handler}
}

/*UpdateJob swagger:route PATCH /v1/jobs/{id} job updateJob

UpdateJob update job API

*/
type UpdateJob struct {
	Context *middleware.Context
	Handler UpdateJobHandler
}

func (o *UpdateJob) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewUpdateJobParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
