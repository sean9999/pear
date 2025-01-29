package pear

import (
	"encoding/json"
	"errors"
	"runtime"
)

// *Pear implements SingleError
var _ SingleError = (*Pear)(nil)

// Pear is a perfect error
type Pear struct {
	child    *Pear
	msg      string
	frame    *StackFrame
	distance int // distance from callstack position
}

// ErrorWithStackFrame is an error with a stack frame
type ErrorWithStackFrame struct {
	Message string `json:"msg"`
	*StackFrame
}

// func (p *Pear) Format(w fmt.State, v rune) {
// 	switch v {
// 	case 'w':
// 		//	do it
// 		io.WriteString(w, p.child.msg)
// 	default:
// 		fmtDirective := fmt.FormatString(w, v)

// 		fmt.Fprintf(w, fmtDirective, p)
// 		//io.WriteString(w, fmtDirective)
// 	}
// }

func (p *Pear) Error() string {
	return p.msg
}

// *Pear implements [errors.Error]
// func (p *Pear) Error() string {
// 	if p.child == nil {
// 		return p.msg
// 	}

// 	//str := strings.Replace(p.msg, "%w", p.child.Error(), -1)

// 	//return fmt.Sprintf(str)

// 	return strings.Replace(p.msg, "%w", p.child.Error(), -1)
// }

func IsPear(e error) bool {
	var pear *Pear
	is := errors.As(e, &pear)

	return is
}

// AsPear casts an error to a *Pear, or does nothing if it already is one
func AsPear(e error, offset int) *Pear {
	if e == nil {
		return nil
	}
	if IsPear(e) {
		return e.(*Pear)
	}
	var pe *Pear
	ok := errors.As(e, &pe)
	if ok {
		pe.ensureStackFrame(offset)

		return pe
	}
	p := Pear{
		child: AsPear(errors.Unwrap(e), offset+1),
		msg:   e.Error(),
	}
	p.Stack(offset)

	return &p
}

// Wrap wraps an error, making it the child of our *Pear
func (p *Pear) Wrap(e error) {
	if e == nil {
		return
	}
	var pe *Pear
	ok := errors.As(e, &pe)
	if ok {
		p.child = pe
	} else {
		p.child = AsPear(e, 1)
	}
	p.ensureStackFrame(0)
}

func (p *Pear) Unwrap() error {
	if p == nil {
		return nil
	}

	return p.child
}

// ensure that a *Pear has a stackFrame if it doesn't already
func (p *Pear) ensureStackFrame(offset int) {
	if p.frame == nil {
		p.Stack(offset)
	}
}

// Unwind unwraps *Pear recursively until there are no more children
func (p *Pear) Unwind() []*Pear {
	pears := make([]*Pear, 0, 8) // 8 seems about right for a default size for stack trace
	q := p
	for q != nil {
		pears = append(pears, q)
		q = q.child
	}

	return pears
}

// Trace returns all the [StackFrame]s from the *Pear, down through its descendants
func (p *Pear) Trace() []ErrorWithStackFrame {
	pears := p.Unwind()
	frames := make([]ErrorWithStackFrame, len(pears))
	for i, p := range pears {
		frames[i] = ErrorWithStackFrame{
			StackFrame: p.frame,
			Message:    p.msg,
		}
	}

	return truncateRecords(frames)
}

// Dump a JSON formatted output of the stack trace
func (p *Pear) Dump() string {
	j, _ := json.MarshalIndent(p.Trace(), "", "\t")

	return string(j)
}

// A StackFrame is a record of program execution.
type StackFrame struct {
	Func           string  `json:"func"`
	File           string  `json:"file"`
	Line           int     `json:"line"`
	ProgramCounter uintptr `json:"-"`
}

// Stack produces (or just returns) a StackFrame for a *Pear
func (p *Pear) Stack(offset int) *StackFrame {
	if p.frame == nil {
		pc, file, line, _ := runtime.Caller(p.distance + offset + 2)
		sf := StackFrame{
			File:           file,
			Line:           line,
			ProgramCounter: pc,
			Func:           runtime.FuncForPC(pc).Name(),
		}
		p.frame = &sf
	}

	return p.frame
}

func (p *Pear) Throw(offset int) *Pear {
	p.ensureStackFrame(offset)

	return p
}

func Dump(e error) string {
	pear := AsPear(e, 0)

	return pear.Dump()
}
