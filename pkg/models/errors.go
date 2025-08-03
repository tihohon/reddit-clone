package models

type NotFoundError struct{}

func (m *NotFoundError) Error() string {
	return "not found"
}

type NoAuthError struct{}

func (m *NoAuthError) Error() string {
	return "no auth"
}

type SignError struct{}

func (m *SignError) Error() string {
	return "failed do sign data"
}

type InvalidValueError struct{}

func (m *InvalidValueError) Error() string {
	return "invalid value"
}
