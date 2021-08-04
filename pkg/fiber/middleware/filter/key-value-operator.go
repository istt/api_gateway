package filter

import (
	"fmt"
	"strings"
)

const (
	// Value Operators
	OpIsNull     = "isnull"     // Is Null
	OpIsNotNull  = "isnotnull"  // Is not Null
	OpIsEmpty    = "isempty"    // Is Empty
	OpIsNotEmpty = "isnotempty" // Is not Empty
	// String Operators
	OpEq             = "eq"             // Is Equal To
	OpNeq            = "neq"            // Not Equals To
	OpStartsWith     = "startswith"     // Starts With
	OpContains       = "contains"       // Contains
	OpEndsWith       = "endswith"       // Ends With
	OpDoesNotContain = "doesnotcontain" // Does Not Contain
	// Numeric Operators
	OpLt  = "lt"  // Less Than
	OpLte = "lte" // Less Than or Equal
	OpGte = "gte" // Greater Than or Equal
	OpGt  = "gt"  // Greater Than
	// Range Operators
	OpIn         = "in"         // In array of values
	OpNotIn      = "notin"      // Not In array of values
	OpBetween    = "between"    // Between values
	OpNotBetween = "notbetween" // Not Between values
)

type KeyvalueOperator struct {
	Key      string
	Operator string
	Value    string
	Values   []string
}

// Unmarshal the Query Parameters
func (kov *KeyvalueOperator) UnmarshalQueryParam(key string, values ...string) error {
	value := values[0]
	kov.Values = values
	if o, k, err := kov.match(key); err == nil { // operator is located in the key part
		kov.Key = k
		kov.Operator = o
		kov.Value = value
	} else if o, v, err := kov.match(value); err == nil { // operator is located in the value part
		kov.Key = key
		kov.Operator = o
		kov.Value = v
	} else {
		kov.Key = key
		kov.Value = value
		kov.Operator = "eq"
	}
	return kov.Validate()
}

// Validate if the KeyvalueOperator set is correct
func (kov KeyvalueOperator) Validate() error {
	switch kov.Operator {
	case OpIsNull, OpIsNotNull, OpIsEmpty, OpIsNotEmpty:
		return nil
	case OpEq, OpNeq, OpStartsWith, OpContains, OpEndsWith, OpDoesNotContain, OpLt, OpLte, OpGte, OpGt:
		if kov.Value == "" {
			return fmt.Errorf("operator %s require exactly 1 parameter", kov.Operator)
		}
	case OpBetween, OpNotBetween:
		if len(kov.Values) != 2 {
			return fmt.Errorf("operator %s require exactly 2 parameters", kov.Operator)
		}
		return nil
	case OpIn, OpNotIn:
		if len(kov.Values) == 0 {
			return fmt.Errorf("operator %s require at least one parameter", kov.Operator)
		}
	default:
		return fmt.Errorf("unknown operator: %s", kov.Operator)
	}
	return nil
}

// extract the value from string, or skip it
func (kov KeyvalueOperator) match(s string) (string, string, error) {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s[i:], ")")
		if j >= 0 {
			return s[0:i], s[i+1 : j+i], nil
		}
	}
	return "", "", fmt.Errorf("not found")
}

func (kov KeyvalueOperator) IsSingle() bool {
	switch kov.Operator {
	case OpEq, OpNeq, OpStartsWith, OpContains, OpEndsWith, OpDoesNotContain, OpLt, OpLte, OpGte, OpGt:
		return true
	default:
		return false
	}
}

func (kov KeyvalueOperator) IsDouble() bool {
	switch kov.Operator {
	case OpBetween, OpNotBetween:
		return true
	default:
		return false
	}
}

func (kov KeyvalueOperator) IsMultiple() bool {
	switch kov.Operator {
	case OpIn, OpNotIn:
		return true
	default:
		return false
	}
}
