package filter

import (
	"errors"
)

// PageRequest model a pageable
// https://docs.spring.io/spring-data/commons/docs/current/api/org/springframework/data/domain/PageRequest.html
type PageRequest struct {
	Page int // zero-based page index, must not be negative.
	Size int // the size of the page to be returned, must be greater than 0.
}

// Validate a page request
func (p PageRequest) Validate() error {
	if p.Page < 0 {
		return errors.New("page should start from 0")
	}
	if p.Size <= 0 {
		return errors.New("size should greater than 1")
	}
	return nil
}
