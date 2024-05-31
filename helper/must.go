package helper

func Must[T any](in T, err error) T {
	PanicIfError(err, "must not return error")
	return in
}
