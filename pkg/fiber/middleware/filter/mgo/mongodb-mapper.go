package mgo

import (
	"fmt"

	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FilterMongo struct {
	criteria *filter.Filter
}

// Constructor of a new FilterSQL
func NewFilterMongo(criteria *filter.Filter) FilterMongo {
	return FilterMongo{criteria: criteria}
}

// MarshalSQL provide the operator for SQL
func (f FilterMongo) MarshalBSON() (bson.M, *options.FindOptions, error) {
	// build where conditions
	filterOps := []bson.M{}
	for _, kov := range (*f.criteria).Filters {
		switch {
		case kov.IsSingle():
			if item, err := f.filterToMongo(kov.Key, kov.Operator, kov.Value); err == nil {
				filterOps = append(filterOps, item)
			}
		case kov.IsDouble():
			if item, err := f.filterMongoDouble(kov.Key, kov.Operator, kov.Values[0], kov.Values[1]); err == nil {
				filterOps = append(filterOps, item)
			}
		case kov.IsMultiple():
			if item, err := f.filterMongoMultiple(kov.Key, kov.Operator, kov.Values[0], kov.Values[1]); err == nil {
				filterOps = append(filterOps, item)
			}
		default:
			if item, err := f.filterToMongo(kov.Key, kov.Operator, kov.Value); err == nil {
				filterOps = append(filterOps, item)
			}
		}

	}
	filterD := bson.M{}
	if len(filterOps) > 0 {
		filterD = bson.M{"$and": filterOps}
	}

	opt := &options.FindOptions{}

	// sorting
	if (*f.criteria).Sort != nil {
		if (*f.criteria).Sort.Direction == filter.SORT_ASC {
			opt.Sort = bson.M{(*f.criteria).Sort.Property: 1}
		} else {
			opt.Sort = bson.M{(*f.criteria).Sort.Property: -1}
		}
	}

	// pagination
	if (*f.criteria).Page != nil {
		size := int64((*f.criteria).Page.Size)
		opt.Limit = &size
		skip := int64((*f.criteria).Page.Size * ((*f.criteria).Page.Page - 1))
		opt.Skip = &skip
	}
	return filterD, opt, nil
}

func (f FilterMongo) filterToMongo(attribute, op, value string) (bson.M, error) {
	switch op {
	case filter.OpIsNull:
		return bson.M{attribute: nil}, nil
	case filter.OpIsNotNull:
		return bson.M{attribute: bson.M{"$ne": nil}}, nil
	case filter.OpIsEmpty:
		return bson.M{attribute: ""}, nil
	case filter.OpIsNotEmpty:
		return bson.M{attribute: bson.M{"$ne": ""}}, nil
	case filter.OpEq:
		return bson.M{attribute: value}, nil
		// return fmt.Sprintf("%s = ?", attribute), , nil
	case filter.OpNeq:
		return bson.M{attribute: bson.M{"$ne": value}}, nil
	case filter.OpStartsWith:
		return bson.M{attribute: primitive.Regex{Pattern: "^" + value, Options: "i"}}, nil
	case filter.OpContains:
		return bson.M{attribute: primitive.Regex{Pattern: "/" + value + "/", Options: "i"}}, nil
	case filter.OpEndsWith:
		return bson.M{attribute: primitive.Regex{Pattern: value + "$", Options: "i"}}, nil
	case filter.OpDoesNotContain:
		return bson.M{attribute: primitive.Regex{Pattern: "^((?!" + value + ").)*$", Options: "i"}}, nil
	case filter.OpLt:
		return bson.M{attribute: bson.M{"$lt": value}}, nil
	case filter.OpLte:
		return bson.M{attribute: bson.M{"$lte": value}}, nil
	case filter.OpGt:
		return bson.M{attribute: bson.M{"$gt": value}}, nil
	case filter.OpGte:
		return bson.M{attribute: bson.M{"$gte": value}}, nil
	default:
		return nil, fmt.Errorf("not a single value operator")
	}
}

// Apply on operator that only contains 2 values
func (f FilterMongo) filterMongoDouble(attribute, op, start, end string) (bson.M, error) {
	switch op {
	case filter.OpBetween:
		return bson.M{attribute: bson.M{"$and": []bson.M{{"$gte": start}, {"$lte": end}}}}, nil
	case filter.OpNotBetween:
		return bson.M{attribute: bson.M{"$or": []bson.M{{"$lt": start}, {"$gt": end}}}}, nil
	default:
		return nil, fmt.Errorf("not a single value operator")
	}
}

// Apply on operator that only contains 2 values
func (f FilterMongo) filterMongoMultiple(attribute, op string, values ...string) (bson.M, error) {
	switch op {
	case filter.OpIn:
		return bson.M{attribute: bson.M{"$in": values}}, nil
	case filter.OpNotIn:
		return bson.M{attribute: bson.M{"$nin": values}}, nil
	default:
		return nil, fmt.Errorf("not a single value operator")
	}
}
