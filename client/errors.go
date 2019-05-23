package client

import "fmt"

type statusCodeError struct {
	StatusCode int
}

func (s *StatusCodeError) Error() string {
	return fmt.Sprintf("status code %d", s.StatusCode)
}

func (s *StatusCode) StatusCode() int {
	return s.StatusCode
}
