package entity

type invalidConfigError struct {
	message string
}

func (e invalidConfigError) Error() string {
	return e.message
}

func ErrInvalidConfig(message string) error {
	return invalidConfigError{message: message}
}
