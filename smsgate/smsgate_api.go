package smsgate

import (
	"net/http"
)

// ErrorCode is application error code
type ErrorCode int

const (
	// RCApplicationError is internal application error.
	RCApplicationError ErrorCode = 101

	// RCEncodingError is error returned when encoding
	// is not supported or message not encoded with given encoding.
	RCEncodingError ErrorCode = 102

	// RCNoAccount is error returned when there is
	// no account with given username/password.
	RCNoAccount ErrorCode = 103

	// RCIPNotAllowed is error returned when sending
	// from clients IP address is not allowed.
	RCIPNotAllowed ErrorCode = 104

	// RCThrottlingError is error returned when there are
	// too many messages submitted withing short period of time.
	RCThrottlingError ErrorCode = 105

	// RCBlacklistedSender is error returned when sender
	// contains words blacklisted on destination.
	RCBlacklistedSender ErrorCode = 106

	// RCInvalidSender is error returned when sender contains
	// illegal characters.
	RCInvalidSender ErrorCode = 107

	// RCMessageTooLong is error returned when message text is
	// too long.
	RCMessageTooLong ErrorCode = 108

	// RCBadContentFormat is error returned when format of
	// text/content parameter is wrong.
	RCBadContentFormat ErrorCode = 109

	// RCMissingMandatoryParameter is error returned when
	// mandatory parameter is missing.
	RCMissingMandatoryParameter ErrorCode = 110

	// RCUnknownMessageType is error returned when message type
	// is unknown.
	RCUnknownMessageType ErrorCode = 111

	// RCBadParameterValue is error returned when parameter has
	// bad value.
	RCBadParameterValue ErrorCode = 112

	// RCNoCredit is error returned when there is no more credit
	// on account balance.
	RCNoCredit ErrorCode = 113

	// RCNoRoute is error returned when there is no route.
	RCNoRoute ErrorCode = 114

	// RCConcatError is error returned when message cannot be
	// split into parts.
	RCConcatError ErrorCode = 115
)

// APIError is error returned by API
type APIError struct {
	message string
	code    ErrorCode
}

// NewAPIError creates new API error
func NewAPIError(message string, code ErrorCode) *APIError {
	return &APIError{message: message, code: code}
}

// Error returns error message (errors interface)
func (err APIError) Error() string {
	return err.message
}

// Code returns API Error Code
func (err *APIError) Code() ErrorCode {
	return err.code
}

// IsAPIError returns true if given error is API error
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// API is interface for accessing SMS BulkHTTP service
type API interface {
	// Send SMS over API
	Send(req *SMSRequest) (*SMSResponse, error)

	// ParseDeliveryReport should be used to transform received DLR
	// with DeliveryReport or error
	ParseDeliveryReport(req *http.Request) (*DeliveryReport, error)
}
