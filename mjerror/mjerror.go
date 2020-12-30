package mjerror

import "github.com/pkg/errors"

var (
	ErrPlyNotFound = errors.New("find ply failed")
)
