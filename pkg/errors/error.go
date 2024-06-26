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

type DuplicateKeyError struct {
	Wrapped error
}

func (d DuplicateKeyError) Error() string {
	return "invalid argument : service name already existss"
}

func (d DuplicateKeyError) Unwrap() error {
	return d.Wrapped
}
