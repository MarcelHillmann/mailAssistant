package errors

type emptyError struct {
}

func (e emptyError) Error() string {
	return "is empty"
}

// NewEmpty is a factory for emptyError
func NewEmpty() error {
	return emptyError{}
}

// IsEmpty is checking if it is a emptyError instance
func IsEmpty(inErr error) bool {
	switch inErr.(type) {
	case *emptyError:
		return true
	case emptyError:
		return true
	default:
		return false
	}
}
