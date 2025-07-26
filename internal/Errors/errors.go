package errors

import "errors"

var (
	ErrCouldNotCreateFile   = errors.New("error creating file")
	ErrCouldNotCreateDialer = errors.New("error creating dialer")
	ErrTopicHasNoMessages   = errors.New("error topic has no messages")
	ErrCACertError          = errors.New("error appending CA certificate")
	ErrRequiredFields       = errors.New("error required fields/s empty")
	ErrUnknownType          = errors.New("error unknown kafka type")
)
