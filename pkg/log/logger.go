package log

import (
	"strings"

	"github.com/go-logr/logr"
)

// logger is a wrapper around a logr.Logger that writes log messages
// about the API discovery and code generation process.
type logger struct {
	log        logr.Logger
	blockDepth int
}

// IsDebugEnabled returns true when the underlying logger is configured to
// write debug messages, false otherwise.
func (l *logger) IsDebugEnabled() bool {
	return l.log.V(1).Enabled()
}

// WithValues adapts the internal logger with a set of additional values
func (l *logger) WithValues(
	values ...interface{},
) {
	l.log = l.log.WithValues(values...)
}

// Debug writes a supplied log message if debug logging is enabled
func (l *logger) Debug(
	msg string,
	vals ...interface{},
) {
	l.log.V(1).Info(msg, vals...)
}

// Info writes a supplied log message about a resource that includes a
// set of standard log values for the resource's kind, namespace, name, etc
func (l *logger) Info(
	msg string,
	vals ...interface{},
) {
	l.log.V(0).Info(msg, vals...)
}

// Enter logs an entry to a function or code block
func (l *logger) Enter(
	name string, // name of the function or code block we're entering
	vals ...interface{},
) {
	if l.log.V(1).Enabled() {
		l.blockDepth++
		depth := strings.Repeat(">", l.blockDepth)
		msg := depth + " " + name
		l.log.V(1).Info(msg, vals...)
	}
}

// Exit logs an exit from a function or code block
func (l *logger) Exit(
	name string, // name of the function or code block we're exiting
	err error,
	vals ...interface{},
) {
	if l.log.V(1).Enabled() {
		depth := strings.Repeat("<", l.blockDepth)
		msg := depth + " " + name
		if err != nil {
			vals = append(vals, "error")
			vals = append(vals, err)
		}
		l.log.V(1).Info(msg, vals...)
		l.blockDepth--
	}
}

// Trace logs an entry to a function or code block and returns a functor
// that can be called to log the exit of the function or code block
func (l *logger) Trace(
	name string,
	vals ...interface{},
) TraceExiter {
	l.Enter(name, vals...)
	f := func(err error, args ...interface{}) {
		l.Exit(name, err, args...)
	}
	return f
}

// New returns a logger that can write log messages about API discovery and
// code generation processes.
func New(
	log logr.Logger,
	vals ...interface{},
) *logger {
	return &logger{
		log:        log.WithValues(vals...),
		blockDepth: 0,
	}
}
