package sql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
)

// FilterSQL help convert FilterDTO into suitable data for GORM library
type FilterSQL struct {
	filter.Filter
	criteria        []string
	sqlClause       string
	namedParameters map[string]interface{}
}

// FilterSQL build the SQL with given params compatible with GORM
func (s *FilterSQL) MarshalSQLWhere() (*string, *map[string]interface{}, error) {
	// + build the where clause
	s.criteria = make([]string, 0)
	s.namedParameters = make(map[string]interface{})
	for _, f := range s.Filters {
		if err := s.marshalOperator(f); err != nil {
			return nil, nil, err
		}
	}
	s.sqlClause = strings.Join(s.criteria, " AND ")
	return &s.sqlClause, &s.namedParameters, nil
}

// MarshalSQLOrderBy add the Order By clause
func (s *FilterSQL) MarshalSQLOrderBy() string {
	if s.Sort == nil {
		return ""
	}
	if s.Sort.Direction == filter.SORT_DESC {
		return fmt.Sprintf("%s %s", s.Sort.Property, filter.SORT_DESC)
	}
	return fmt.Sprintf("%s %s", s.Sort.Property, filter.SORT_ASC)
}

func (s *FilterSQL) marshalOperator(kov *filter.KeyvalueOperator) error {
	switch kov.Operator {
	case filter.OpIsNull: // Is Null
		s.criteria = append(s.criteria, fmt.Sprintf("%s IS NULL", kov.Key))
	case filter.OpIsNotNull: // Is not Null
		s.criteria = append(s.criteria, fmt.Sprintf("%s IS NOT NULL", kov.Key))
	case filter.OpIsEmpty: // Is Empty
		s.criteria = append(s.criteria, fmt.Sprintf("%s = ''", kov.Key))
	case filter.OpIsNotEmpty: // Is not Empty
		s.criteria = append(s.criteria, fmt.Sprintf("%s != ''", kov.Key))
	case filter.OpEq: // Is Equal To
		s.namedParameters[kov.Key] = kov.Value
		s.criteria = append(s.criteria, fmt.Sprintf("%s = @%s", kov.Key, kov.Key))
	case filter.OpNeq: // Not Equals To
		s.namedParameters[kov.Key] = kov.Value
		s.criteria = append(s.criteria, fmt.Sprintf("%s != @%s", kov.Key, kov.Key))
	case filter.OpStartsWith: // Starts With
		s.namedParameters[kov.Key] = kov.Value + "%"
		s.criteria = append(s.criteria, fmt.Sprintf("%s LIKE @%s", kov.Key, kov.Key))
	case filter.OpContains: // Contains
		s.namedParameters[kov.Key] = "%" + kov.Value + "%"
		s.criteria = append(s.criteria, fmt.Sprintf("%s LIKE @%s", kov.Key, kov.Key))
	case filter.OpEndsWith: // Ends With
		s.namedParameters[kov.Key] = "%" + kov.Value
		s.criteria = append(s.criteria, fmt.Sprintf("%s LIKE @%s", kov.Key, kov.Key))
	case filter.OpDoesNotContain: // Does Not Contain
		s.namedParameters[kov.Key] = "%" + kov.Value + "%"
		s.criteria = append(s.criteria, fmt.Sprintf("%s NOT LIKE @%s", kov.Key, kov.Key))
	case filter.OpLt: // Less Than
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.namedParameters[kov.Key] = num
		s.criteria = append(s.criteria, fmt.Sprintf("%s < @%s", kov.Key, kov.Key))
	case filter.OpLte: // Less Than or Equal
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.namedParameters[kov.Key] = num
		s.criteria = append(s.criteria, fmt.Sprintf("%s <= @%s", kov.Key, kov.Key))
	case filter.OpGte: // Greater Than or Equal
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.namedParameters[kov.Key] = num
		s.criteria = append(s.criteria, fmt.Sprintf("%s >= @%s", kov.Key, kov.Key))
	case filter.OpGt: // Greater Than
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.namedParameters[kov.Key] = num
		s.criteria = append(s.criteria, fmt.Sprintf("%s > @%s", kov.Key, kov.Key))

	case filter.OpIn: // In array of values
		if len(kov.Values) == 0 {
			return fmt.Errorf("the IN operator need a list of values")
		}
		s.namedParameters[kov.Key] = kov.Values
		s.criteria = append(s.criteria, fmt.Sprintf("%s IN @%s", kov.Key, kov.Key))
	case filter.OpNotIn: // Not In array of values
		if len(kov.Values) == 0 {
			return fmt.Errorf("the IN operator need a list of values")
		}
		s.namedParameters[kov.Key] = kov.Values
		s.criteria = append(s.criteria, fmt.Sprintf("%s NOT IN @%s", kov.Key, kov.Key))
	case filter.OpBetween: // Between values
		s.namedParameters[kov.Key+"Start"] = kov.Values[0]
		s.namedParameters[kov.Key+"End"] = kov.Values[1]
		s.criteria = append(s.criteria, fmt.Sprintf("%s BETWEEN @%sStart AND @%sEnd", kov.Key, kov.Key, kov.Key))
	case filter.OpNotBetween: // Not Between values
		s.namedParameters[kov.Key+"Start"] = kov.Values[0]
		s.namedParameters[kov.Key+"End"] = kov.Values[1]
		s.criteria = append(s.criteria, fmt.Sprintf("%s NOT BETWEEN @%sStart AND @%sEnd", kov.Key, kov.Key, kov.Key))
	}
	return nil
}
