package pear

// an Error is an error with a stack trace.
type Error interface {
	Error() string
	Wrap(error)
	Unwind() []*Pear
	Trace() []ErrorWithStackFrame
	Dump() string
}

// a SingleError is the most common type of error
type SingleError interface {
	Error
	Unwrap() error
	Stack(int) *StackFrame
}

// a MultiError is an error made up of a list of errors
type MultiError interface {
	Error
	Unwrap() []error
}
