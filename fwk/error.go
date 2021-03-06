package fwk

import (
	"go-hep.org/x/hep/fwk/utils/errstack"
)

// Errorf formats according to a format specifier and returns the string as
// a value that satisfies error, together with the associated stack trace.
func Errorf(format string, args ...interface{}) error {
	return errstack.Newf(format, args...)
}

// Error returns the original error with the associated stack trace.
func Error(err error) error {
	return errstack.New(err)
}

// EOF
