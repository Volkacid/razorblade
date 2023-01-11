package storage

import (
	"fmt"
)

type NFError struct {
	Err string
}

func (nfe *NFError) Error() string {
	return fmt.Sprint(nfe.Err)
}
func NotFoundError() error {
	return &NFError{Err: "storage: not found"}
}

type DuplicateError struct {
	Err string
}

func (dupErr *DuplicateError) Error() string {
	return fmt.Sprint(dupErr.Err)
}
func FoundDuplicateError() error {
	return &DuplicateError{Err: "storage: short URL for this value already exists"}
}

type DeletedError struct {
	Err string
}

func (delErr *DeletedError) Error() string {
	return fmt.Sprint(delErr.Err)
}

func ValueDeletedError() error {
	return &DeletedError{Err: "storage: this value is deleted"}
}
