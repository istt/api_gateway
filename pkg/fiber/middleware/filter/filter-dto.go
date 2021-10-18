package filter

type Filter struct {
	Filters []*KeyvalueOperator
	Sort    *Sort
	Page    *PageRequest
}

// TODO: Find a way to build query based on above listing
