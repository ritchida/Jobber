package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-swagger/go-swagger/strfmt"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/validate"
)

/*JobStatus job status

swagger:model JobStatus
*/
type JobStatus string

// for schema
var jobStatusEnum []interface{}

func (m JobStatus) validateJobStatusEnum(path, location string, value JobStatus) error {
	if jobStatusEnum == nil {
		var res []JobStatus
		if err := json.Unmarshal([]byte(`["created","running","succeeded","failed","unknown"]`), &res); err != nil {
			return err
		}
		for _, v := range res {
			jobStatusEnum = append(jobStatusEnum, v)
		}
	}
	if err := validate.Enum(path, location, value, jobStatusEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this job status
func (m JobStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJobStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
