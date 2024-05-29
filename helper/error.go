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
