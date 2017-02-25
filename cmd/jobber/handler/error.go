package handler

import (
	"fmt"

	apimodel "github.com/ritchida/jobber/generated/jobber/models"
)

func createServiceError(err error, statusCode int64) apimodel.Error {
	se := apimodel.Error{
		Code:    statusCode,
		Message: fmt.Sprintf("%v", err),
	}
	return se
}
