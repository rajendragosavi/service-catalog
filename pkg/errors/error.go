package errors

type ErrorArgument struct {
	Wrapped error
}
type DuplicateKeyError struct {
	Wrapped error
}
type ObjectNotFoundError struct {
	Wrapped error
}

type SystemError struct {
	Wrapped error
}

func (e ErrorArgument) Error() string {
	return "invalid argument"
}

func (d DuplicateKeyError) Error() string {
	return "invalid argument : service name already exists"
}

func (o ObjectNotFoundError) Error() string {
	return "invalid argument : provided service name does not exist"
}

func (s SystemError) Error() string {
	return "service catalog is down. please retry after some time, if issue still persists reach out to service catalog admin team"
}

func (s SystemError) Unwrap() error {
	return s.Wrapped
}
