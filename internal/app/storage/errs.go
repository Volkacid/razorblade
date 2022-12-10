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

func (de *DuplicateError) Error() string {
	return fmt.Sprint(de.Err)
}
func FoundDuplicateError() error {
	return &DuplicateError{Err: "storage: short URL for this value already exists"}
}
