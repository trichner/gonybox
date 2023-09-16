package errwrap

import "strings"

func Wrap(msg string, wrapped error) error {
	return &wrapError{
		msg: msg,
		err: wrapped,
	}
}

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	if e.err == nil {
		return e.msg
	}

	var buf strings.Builder
	buf.WriteString(e.msg)
	buf.WriteString(": ")
	buf.WriteString(e.err.Error())
	return buf.String()
}

func (e *wrapError) Unwrap() error {
	return e.err
}
