package pear

import (
	"errors"
	"fmt"
)

/**
 *	Methods to assist in using this package as a drop-in replacement for errors
 **/

// drop-in replacement for [errors.New]
func New(msg string) *Pear {
	p := Defer(msg)
	p.ensureStackFrame(1)
	return p
}

// use this for sentinal errors.
// ex: var MyError = errors.New("something that is thrown later")
// ex: var MyError = pear.Defer("something that is thrown later")
func Defer(msg string) *Pear {
	p := &Pear{msg: msg, distance: 1}
	return p
}

// drop-in replacement for [errors.As]
func As(err error, target any) bool {
	return errors.As(err, target)
}

// drop-in replacement for [errors.Is]
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// drop-in replacement for [errors.Join]
func Join(errs ...error) error {
	pears := make([]*Pear, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			pears = append(pears, AsPear(e, 0))
		}
	}
	if len(pears) == 0 {
		return nil
	}
	return &Multipear{pears}
}

//nolint:err113
func Errorf(fmtstring string, args ...any) error {
	var n int
	for i, a := range args {
		e, isError := a.(error)
		if isError && e != nil {
			n++
			p := AsPear(e, n)
			args[i] = p
		}
	}

	return fmt.Errorf(fmtstring, args...)
}
