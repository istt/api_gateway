package sql

import (
	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
	"gorm.io/gorm"
)

type FilterSQL struct {
	criteria *filter.Filter
	db       *gorm.DB
}

// Constructor of a new FilterSQL
func NewFilterSQL(db *gorm.DB, criteria *filter.Filter) FilterSQL {
	return FilterSQL{criteria: criteria, db: db}
}

// MarshalSQL provide the operator for SQL
func (f FilterSQL) MarshalSQL() error {
	// build where conditions
	for _, kov := range (*f.criteria).Filters {
		switch kov.Operator {
		case "isnull":
			f.db.Where(kov.Key + " IS NULL")
		case "isnotnull":
			f.db.Where(kov.Key + " IS NOT NULL")
		case "isempty":
			f.db.Where(kov.Key + " = ''")
		case "isnotempty":
			f.db.Where(kov.Key + " <> ''")
		case "eq":
			f.db.Where(kov.Key+" = ?", kov.Value)
		case "neq":
			f.db.Where(kov.Key+" <> ?", kov.Value)
		case "startswith":
			f.db.Where(kov.Key+" LIKE ?", "%"+kov.Value)
		case "contains":
			f.db.Where(kov.Key+" LIKE ?", "%"+kov.Key+"%")
		case "endswith":
			f.db.Where(kov.Key+" LIKE ?", kov.Value+"%")
		case "doesnotcontain":
			f.db.Where(kov.Key+" NOT LIKE ?", "%"+kov.Key+"%")
		case "lt":
			f.db.Where(kov.Key+" < ?", kov.Value)
		case "lte":
			f.db.Where(kov.Key+" <= ?", kov.Value)
		case "gt":
			f.db.Where(kov.Key+" > ?", kov.Value)
		case "gte":
			f.db.Where(kov.Key+" >= ?", kov.Value)
		case "in":
			f.db.Where(kov.Key+" IN ?", kov.Value)
		case "notin":
			f.db.Where(kov.Key+" NOT IN ?", kov.Value)
		case "between":
			f.db.Where(kov.Key+" BETWEEN ? AND ?", kov.Value)
		case "notbetween":
			f.db.Where(kov.Key+" NOT BETWEEN ? AND ?", kov.Value)
		}
	}

	// sorting
	if (*f.criteria).Sort != nil {
		if orderClause, err := (*f.criteria).Sort.MarshalText(); err == nil {
			f.db.Order(orderClause)
		}
	}

	// pagination
	if (*f.criteria).Page != nil {
		f.db.Limit((*f.criteria).Page.Size)
		f.db.Offset((*f.criteria).Page.Size * ((*f.criteria).Page.Page - 1))
	}
	return nil
}
