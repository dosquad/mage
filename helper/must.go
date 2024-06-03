package helper

// Must takes a function that returns T, error and if the error is not nil
// panics, otherwise returns just T.
func Must[T any](in T, err error) T {
	PanicIfError(err, "must not return error")
	return in
}
