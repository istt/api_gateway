package filter

import (
	"fmt"
	"strings"
)

const (
	SORT_ASC  = "ASC"
	SORT_DESC = "DESC"
)

type Sort struct {
	Property  string
	Direction string
}

// UnmarshalText convert `property,direction` into property, direction attribute of sort
func (s *Sort) UnmarshalText(text []byte) error {
	txt := strings.TrimSpace(string(text))
	if len(txt) == 0 {
		return fmt.Errorf("empty text")
	}
	comma := strings.Index(txt, ",")
	if comma == -1 {
		s.Property = txt
		s.Direction = SORT_ASC
	} else {
		s.Property = txt[0:comma]
		s.Direction = txt[comma+1:]
	}
	return nil
}

// MarshalText convert the Sort back into text string like `properties,sortDirection`
func (s Sort) MarshalText() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return []byte{}, err
	}
	return []byte(fmt.Sprintf("%s,%s", s.Property, s.Direction)), nil
}

// Validate ensure that the sort contains a valid property and direction
func (s Sort) Validate() error {
	if len(s.Property) == 0 {
		return fmt.Errorf("missing sort property")
	}
	switch strings.ToUpper(s.Direction) {
	case SORT_DESC:
		s.Direction = SORT_DESC
	default:
		s.Direction = SORT_ASC
	}
	return nil
}
