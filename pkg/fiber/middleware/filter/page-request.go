package filter

import (
	"errors"
)

type PageRequest struct {
	Page int
	Size int
}

// Validate a page request
func (p PageRequest) Validate() error {
	if p.Page < 1 {
		return errors.New("page should start from 1")
	}
	if p.Size <= 0 {
		return errors.New("size should greater than 1")
	}
	return nil
}
