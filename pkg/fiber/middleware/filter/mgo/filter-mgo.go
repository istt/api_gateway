package mgo

import (
	"fmt"
	"strconv"

	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FilterBSON help convert FilterDTO into suitable data for GORM library
type FilterBSON struct {
	filter.Filter
	condition bson.D
}

// FilterBSON build the BSON with given params compatible with GORM
func (s *FilterBSON) MarshalBSONCondition() (*bson.D, error) {
	s.condition = bson.D{}
	for _, f := range s.Filters {
		if err := s.marshalOperator(f); err != nil {
			return nil, err
		}
	}

	return &s.condition, nil
}

// MarshalBSONOrderBy add the Order By clause
func (s *FilterBSON) MarshalBSONOrderBy() bson.D {
	if s.Sort == nil {
		return bson.D{}
	}
	if s.Sort.Direction == filter.SORT_DESC {
		return bson.D{{Key: s.Sort.Property, Value: -1}}
	}
	return bson.D{{Key: s.Sort.Property, Value: 1}}
}

func (s *FilterBSON) marshalOperator(kov *filter.KeyvalueOperator) error {
	switch kov.Operator {
	case filter.OpIsNull: // Is Null
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.TypeNull})
	case filter.OpIsNotNull: // Is not Null
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$ne", Value: bson.TypeNull}})
	case filter.OpIsEmpty: // Is Empty
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$exists", Value: false}})
	case filter.OpIsNotEmpty: // Is not Empty
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$exists", Value: true}})
	case filter.OpEq: // Is Equal To
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$eq", Value: kov.Value}})
	case filter.OpNeq: // Not Equals To
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$ne", Value: kov.Value}})
	case filter.OpStartsWith: // Starts With
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: primitive.Regex{Pattern: "^" + kov.Value, Options: "i"}})
	case filter.OpContains: // Contains
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: primitive.Regex{Pattern: kov.Value, Options: "i"}})
	case filter.OpEndsWith: // Ends With
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: primitive.Regex{Pattern: kov.Value + "$", Options: "i"}})
	case filter.OpDoesNotContain: // Does Not Contain
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$not", Value: bson.E{Key: "$regex", Value: primitive.Regex{Pattern: kov.Value, Options: "i"}}}})
	case filter.OpLt: // Less Than
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$lt", Value: num}})
	case filter.OpLte: // Less Than or Equal
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$lte", Value: num}})
	case filter.OpGte: // Greater Than or Equal
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$gte", Value: num}})
	case filter.OpGt: // Greater Than
		num, err := strconv.Atoi(kov.Value)
		if err != nil {
			return err
		}
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$gt", Value: num}})

	case filter.OpIn: // In array of values
		if len(kov.Values) == 0 {
			return fmt.Errorf("the IN operator need a list of values")
		}
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$in", Value: kov.Values}})
	case filter.OpNotIn: // Not In array of values
		if len(kov.Values) == 0 {
			return fmt.Errorf("the IN operator need a list of values")
		}
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$nin", Value: kov.Values}})
	case filter.OpBetween: // Between values
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$and", Value: bson.M{"$gt": kov.Values[0], "$lt": kov.Values[1]}}})
	case filter.OpNotBetween: // Not Between values
		s.condition = append(s.condition, bson.E{Key: kov.Key, Value: bson.E{Key: "$not", Value: bson.M{"$gt": kov.Values[0], "$lt": kov.Values[1]}}})
	}
	return nil
}
