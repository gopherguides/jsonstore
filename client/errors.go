package client

import "fmt"

type statusCodeError struct {
	Code int
}

func (s *statusCodeError) Error() string {
	return fmt.Sprintf("status code %d", s.Code)
}

func (s *statusCodeError) StatusCode() int {
	return s.Code
}
