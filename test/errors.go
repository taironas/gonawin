package gonawintest

// ErrString returns the string representation of an error.
func ErrString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
