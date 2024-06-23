package errors

type ErrorArgument struct {
	Wrapped error
}

func (e ErrorArgument) Error() string {
	return "invalid argument"
}

func (e ErrorArgument) Unwrap() error {
	return e.Wrapped
}
