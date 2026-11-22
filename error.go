package errlay

import (
	"fmt"
)

const (
	message = "|layer %s - %d| %s"
)

func New(layer string, code uint8, msg string, err error) error {
	return &Error{
		layer: layer,
		code:  code,
		msg:   msg,
		err:   err,
	}
}

type Error struct {
	layer string
	code  uint8
	msg   string
	err   error
}

func (e *Error) Error() string {
	msg := fmt.Sprintf(message, e.layer, e.code, e.msg)
	if e.err != nil {
		msg += " " + e.err.Error()
	}
	return msg
}

/* Expose Methods */
func (e *Error) Layer() string {
	return e.layer
}

func (e *Error) Code() uint8 {
	return e.code
}

func (e *Error) Message() string {
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.err
}

/* Is Chain */
func (e *Error) Is(target error) bool {
	t, ok := target.(*errTarget)
	if !ok {
		return false
	}

	isLayer := t.layer == nil || *t.layer == e.layer
	isCode := t.code == nil || *t.code == e.code
	return isLayer && isCode
}
