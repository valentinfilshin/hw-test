package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	fmt.Println(val)
	// Place your code here.
	return nil
}
