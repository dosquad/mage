package helper

import (
	"fmt"

	"github.com/fatih/color"
)

// PanicIfError panics if passed error is not nil.
//
// prependStr is optional string to prepend to panic error.
func PanicIfError(err error, prependStr string) {
	if err != nil {
		es := color.RedString("%s", err.Error())
		if prependStr != "" {
			panic(fmt.Errorf("%s : %s", color.RedString("%s", prependStr), es))
		}
		panic(es)
	}
}

func IfErrorf(format string, a ...any) error {
	for _, i := range a {
		if v, ok := i.(error); ok && v != nil {
			return fmt.Errorf(format, a...)
		}
	}

	return nil
}

func IfErrorfWrap[T any](format string) func(T, error) (T, error) {
	return func(t T, err error) (T, error) {
		return t, IfErrorf(format, err)
	}
}
