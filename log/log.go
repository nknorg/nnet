package log

import logging "github.com/op/go-logging"

// Logger is the logger interface
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// The global logger object
var logger Logger = logging.MustGetLogger("nnet")

// SetLogger sets the global logger object
func SetLogger(l Logger) error {
	logger = l
	return nil
}

// Debug logs to the DEBUG log. Arguments are handled in the manner of fmt.Print.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf logs to the DEBUG log. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info logs to the INFO log. Arguments are handled in the manner of fmt.Print.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logs to the INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warning logs to the WARN log. Arguments are handled in the manner of fmt.Print.
func Warning(args ...interface{}) {
	logger.Warning(args...)
}

// Warningf logs to the WARN log. Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

// Error logs to the ERROR log. Arguments are handled in the manner of fmt.Print.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf logs to the ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
