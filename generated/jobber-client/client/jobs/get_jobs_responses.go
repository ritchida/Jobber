package jobs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-swagger/go-swagger/client"
	"github.com/go-swagger/go-swagger/httpkit"

	strfmt "github.com/go-swagger/go-swagger/strfmt"

	"github.com/ritchida/jobber/generated/jobber-client/models"
)

// GetJobsReader is a Reader for the GetJobs structure.
type GetJobsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *GetJobsReader) ReadResponse(response client.Response, consumer httpkit.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetJobsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetJobsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewGetJobsOK creates a GetJobsOK with default headers values
func NewGetJobsOK() *GetJobsOK {
	return &GetJobsOK{}
}

/*GetJobsOK handles this case with default header values.

list of jobs
*/
type GetJobsOK struct {
	Payload []*models.Job
}

func (o *GetJobsOK) Error() string {
	return fmt.Sprintf("[GET /v1/jobs][%d] getJobsOK  %+v", 200, o.Payload)
}

func (o *GetJobsOK) readResponse(response client.Response, consumer httpkit.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetJobsDefault creates a GetJobsDefault with default headers values
func NewGetJobsDefault(code int) *GetJobsDefault {
	return &GetJobsDefault{
		_statusCode: code,
	}
}

/*GetJobsDefault handles this case with default header values.

Error
*/
type GetJobsDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get jobs default response
func (o *GetJobsDefault) Code() int {
	return o._statusCode
}

func (o *GetJobsDefault) Error() string {
	return fmt.Sprintf("[GET /v1/jobs][%d] getJobs default  %+v", o._statusCode, o.Payload)
}

func (o *GetJobsDefault) readResponse(response client.Response, consumer httpkit.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
