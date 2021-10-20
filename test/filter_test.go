package test

import (
	"testing"

	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
)

func TestFilterMiddleware(t *testing.T) {
	sortStr := "id,asc"
	var sort filter.Sort
	if err := sort.UnmarshalText([]byte(sortStr)); err != nil {
		t.Fatal(err)
	}

	sortStr = "ordered, desc"
	if err := sort.UnmarshalText([]byte(sortStr)); err != nil {
		t.Fatal(err)
	}

	sortStr = "id,ASC"
	if err := sort.UnmarshalText([]byte(sortStr)); err != nil {
		t.Fatal(err)
	}

	sortStr = "id,Desc"
	if err := sort.UnmarshalText([]byte(sortStr)); err != nil {
		t.Fatal(err)
	}

	sortStr = "  "
	if err := sort.UnmarshalText([]byte(sortStr)); err == nil {
		t.Fatal("unable to handle empty strings")
	}

	sortStr = " anonymous  "
	if err := sort.UnmarshalText([]byte(sortStr)); err != nil {
		t.Fatal(err)
	}
}
