package errlay

import "errors"

func Target(layer *string, code *uint8) (target error, err error) {
	if layer == nil && code == nil {
		return nil, errors.New("target must have at least one matching field")
	}
	return &errTarget{
		layer: layer,
		code:  code,
	}, nil
}

type errTarget struct {
	layer *string
	code  *uint8
}

func (t *errTarget) Error() string {
	return ""
}
