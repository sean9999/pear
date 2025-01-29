package pear

import (
	"errors"
	"fmt"
	"regexp"
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

// message from fmtstring, without distractions
func withoutTokens(instr string) string {
	pat := regexp.MustCompile(`\s*\%w\s*`)
	outstr := pat.ReplaceAllString(instr, "")
	return outstr
}

func Errorf(fmtstring string, args ...any) error {
	var n int
	for i, a := range args {
		e, isError := a.(error)
		if isError && e != nil {
			n++
			p := AsPear(e, n)
			args[i] = p
			//thisErr.Wrap(p)
			//thisErr = p
		}
	}
	return fmt.Errorf(fmtstring, args...)
}

// drop-in replacement for [fmt.Errorf]
// func Errorf(fmtstring string, args ...any) Error {
// 	var n int

// 	//err := Defer(withoutTokens(fmtstring))

// 	err := Defer(fmtstring)

// 	thisErr := err
// 	for i, a := range args {
// 		e, isError := a.(error)
// 		if isError && e != nil {
// 			n++
// 			childPear := AsPear(e, n)
// 			args[i] = childPear
// 			thisErr.Wrap(childPear)
// 			thisErr = childPear
// 		}
// 	}
// 	pear := AsPear(err, 1)
// 	if n <= 1 {
// 		return pear
// 	}
// 	return NewMultiPear(pear.Unwind()...)
// }
