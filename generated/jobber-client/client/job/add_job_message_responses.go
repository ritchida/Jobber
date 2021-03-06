package job

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

// AddJobMessageReader is a Reader for the AddJobMessage structure.
type AddJobMessageReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *AddJobMessageReader) ReadResponse(response client.Response, consumer httpkit.Consumer) (interface{}, error) {
	switch response.Code() {

	case 202:
		result := NewAddJobMessageAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewAddJobMessageDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewAddJobMessageAccepted creates a AddJobMessageAccepted with default headers values
func NewAddJobMessageAccepted() *AddJobMessageAccepted {
	return &AddJobMessageAccepted{}
}

/*AddJobMessageAccepted handles this case with default header values.

Accepted
*/
type AddJobMessageAccepted struct {
}

func (o *AddJobMessageAccepted) Error() string {
	return fmt.Sprintf("[POST /v1/jobs/{id}/messages][%d] addJobMessageAccepted ", 202)
}

func (o *AddJobMessageAccepted) readResponse(response client.Response, consumer httpkit.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewAddJobMessageDefault creates a AddJobMessageDefault with default headers values
func NewAddJobMessageDefault(code int) *AddJobMessageDefault {
	return &AddJobMessageDefault{
		_statusCode: code,
	}
}

/*AddJobMessageDefault handles this case with default header values.

Error
*/
type AddJobMessageDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the add job message default response
func (o *AddJobMessageDefault) Code() int {
	return o._statusCode
}

func (o *AddJobMessageDefault) Error() string {
	return fmt.Sprintf("[POST /v1/jobs/{id}/messages][%d] addJobMessage default  %+v", o._statusCode, o.Payload)
}

func (o *AddJobMessageDefault) readResponse(response client.Response, consumer httpkit.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
