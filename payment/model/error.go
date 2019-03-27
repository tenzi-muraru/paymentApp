package model

import (
	"fmt"
)

// PaymentError - error response to be returned in case of errors are encountered while handling the HTTP requests
type PaymentError struct {
	Code        string `json:"errorCode"`
	Description string `json:"errorDescription"`
}

// NewBadRequestError - creates a PaymentError for bad requests
func NewBadRequestError(errDescription string) PaymentError {
	return PaymentError{
		Code:        "INVALID_INPUT",
		Description: fmt.Sprintf("Cannot parse request: %v", errDescription),
	}
}

// PaymentNotFound - error response when no payment was found
var PaymentNotFound = PaymentError{
	Code:        "PAYMENT_NOT_FOUND",
	Description: "Cannot find payment with provided ID",
}

// GenericError - generic error response
var GenericError = PaymentError{
	Code:        "GENERIC_ERROR",
	Description: "Failed to process request due to internal error",
}
