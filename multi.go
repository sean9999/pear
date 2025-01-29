package pear

import (
	"encoding/json"
)

// Multipear is here to support [errors.Join]
type Multipear struct {
	errors []*Pear
}

func (m *Multipear) As(target any) bool {
	for _, err := range m.errors {
		if As(err, target) {
			return true
		}
	}

	return false
}

func (m *Multipear) Is(target error) bool {
	for _, err := range m.errors {
		if Is(target, err) {
			return true
		}
	}

	return false
}

func NewMultiPear(errs ...*Pear) *Multipear {
	return &Multipear{errs}
}

func (m *Multipear) Wrap(e error) {
	m.errors = append(m.errors, AsPear(e, 0))
}

func (m *Multipear) Error() string {
	var str string
	for _, e := range m.errors {
		str += e.Error() + "\n"
	}

	return str
}

func (m *Multipear) Unwind() []*Pear {

	return m.errors
}

func (m *Multipear) Unwrap() []error {
	errs := make([]error, len(m.errors))
	for i, merr := range m.Unwind() {
		errs[i] = merr
	}

	return errs
}

func (m *Multipear) Trace() []ErrorWithStackFrame {
	frames := make([]ErrorWithStackFrame, len(m.errors))

	for i, e := range m.errors {
		frames[i] = ErrorWithStackFrame{
			StackFrame: e.frame,
			Message:    e.msg,
		}
	}

	return frames
}

func (m *Multipear) Dump() string {
	j, _ := json.MarshalIndent(m.errors, "", "\t")

	return string(j)
}
