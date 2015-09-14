package gonawintest

// ErrorString returns the string representation of an error.
func ErrorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
