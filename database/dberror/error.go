package dberror

import (
	"fmt"
)

type FindDocError struct {
	err    string
	filter string
}

func (f *FindDocError) Error() string {
	return fmt.Sprintf("not find document, filter(%s), %s", f.filter, f.err)
}

func NewFindDocError(err error, filter interface{}) *FindDocError {
	return &FindDocError{
		err:    err.Error(),
		filter: fmt.Sprintf("%+v", filter),
	}
}
