package scriptparser

import "fmt"

type (
	Severity int

	ValidationError struct {
		Severity Severity
		Line     int
		Column   int
		Message  string
	}
)

const (
	SeverityWarning Severity = iota
	SeverityError
)

func (e ValidationError) String() string {
	level := "warning"
	if e.Severity == SeverityError {
		level = "error"
	}
	return fmt.Sprintf("%s line %d:%d: %s", level, e.Line, e.Column, e.Message)
}
