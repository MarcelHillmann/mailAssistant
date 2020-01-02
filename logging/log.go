package logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
)

// NewLogger is a factory for a new log instance with an autodetected log name
func NewLogger() *Logger {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	return NewNamedLogger(details.Name())
}

// NewNamedLogger is a factory for a new log instance with an given name
func NewNamedLogger(name string) *Logger {
	return newLogger(normalize(name))
}

// NewGlobalLogger is a factory for a new log instance without name and parent
func NewGlobalLogger() *Logger {
	return NewNamedLogger("global")
}

// NewNamedLogWriter is a factory for a new log instance with given name as io.Writer
func NewNamedLogWriter(name string) io.Writer {
	return logWriter{NewNamedLogger(name)}
}

type logWriter struct {
	log *Logger
}

// Writer is converting the byte array to string and delegate to the Debug method
func (l logWriter) Write(p []byte) (n int, err error) {
	msg := "\n" + string(p)
	l.log.Debug(strings.TrimSuffix(msg, "\n"))
	return len(p), nil
}

var loggerRegistry map[string]logLevel

func init(){
	SetLevel("*","")
}

// SetLevel is a key, value pair to set a logLevel programmatically
func SetLevel(name, level string) {
	if name == "*" {
		loggerRegistry = make(map[string]logLevel, 0)
	}else{
		loggerRegistry[name] = stringToLogLevel(level)
	}
}

func newLogger(name string) *Logger {
	return &Logger{name}
}

// Logger represents a Log Entity
type Logger struct {
	name string
}

// Name returns the name of the current entity
func (l Logger) Name() string {
	return l.name
}

// Severe is writing a severe message
func (l Logger) Severe(msg ...interface{}) {
	if l.IsLogLevel(severe) {
		logger(l.name, "SEVERE ", msg)
	}
}

// Warn is writing a WARNING message
func (l Logger) Warn(msg ...interface{}) {
	if l.IsLogLevel(warn) {
		logger(l.name, "WARNING", msg)
	}
}

// Info is writing a INFO message
func (l Logger) Info(msg ...interface{}) {
	if l.IsLogLevel(info) {
		logger(l.name, "INFO   ", msg)
	}
}

// Debug is writing a debug message
func (l Logger) Debug(msg ...interface{}) {
	if l.IsLogLevel(debug) {
		logger(l.name, "DEBUG  ", msg)
	}
}

// Enter is writing a special debug message that represents to enter a method
func (l Logger) Enter() {
	l.Debug(">>")
}

// Leave is writing a special debug message that represents to leave a method
func (l Logger) Leave() {
	l.Debug("<<")
}

// IsLogLevel check the loggerRegistry for a given log level
func (l Logger) IsLogLevel(level logLevel) bool {
	if lvl, ok := loggerRegistry[l.name]; ok && lvl != notExists {
		return lvl <= level
	}
	parents := strings.Split(l.name, ".")
	for i := len(parents) - 1; i > 0; i-- {
		if lvl, ok := loggerRegistry[buildParent(parents[0:i])]; ok && lvl != notExists {
			return lvl <= level
		}
	}
	if lvl, ok := loggerRegistry["global"]; ok && lvl != notExists {
		return lvl <= level
	}
	return false
}

// Panic is writing a severe message and calls a native panic
func (l Logger) Panic(msg ...interface{}) {
	l.Severe(msg...)

	errorMsg := ""
	for _, m := range msg {
		if errorMsg != "" {
			errorMsg += " "
		}
		errorMsg += fmt.Sprint(m)
	}
	panic(errors.New(errorMsg))
}

func buildParent(parents []string) string {
	name := ""
	for _, parent := range parents {
		if name != "" {
			name += "."
		}
		name += parent
	}

	return name
}

func normalize(name string) string {
	var modifyLogic = []func(string) string{
		func(ugly string) string {
			return strings.Replace(ugly, "/", ".", -1)
		}, func(ugly string) string {
			return strings.TrimPrefix(ugly, "mailassistant.logging.testlogger.")
		}, func(ugly string) string {
			return strings.Replace(ugly, "(", "", -1)
		}, func(ugly string) string {
			return strings.Replace(ugly, ")", "", -1)
		}, func(ugly string) string {
			return strings.Replace(ugly, ".*", ".", -1)
		}, func(ugly string) string {
			return strings.TrimSuffix(ugly, ".getLogger")
		}, func(ugly string) string {
			return strings.TrimSuffix(ugly, ".logger")
		}, func(ugly string) string {
			return strings.Replace(ugly, "${project}", "mailAssistant", -1)
		},func(ugly string ) string {
			return strings.TrimSuffix(ugly,"mailassistant.actions.")
		},func(ugly string ) string {
			return strings.TrimSuffix(ugly,"mailAssistant.actions.")
		},

	}

	makeItNicer := name
	for _, modify := range modifyLogic {
		makeItNicer = modify(makeItNicer)
	}

	return makeItNicer
}

func logger(name, level string, msg []interface{}) {
	methodNameUgly := ""
	for i := 2; i < 4; i++ {
		pc, _, _, _ := runtime.Caller(i)
		method := runtime.FuncForPC(pc)
		methodNameUgly = strings.ToLower(method.Name())
		if !strings.HasSuffix(methodNameUgly, "/logging.logger.panic") &&
			!strings.HasSuffix(methodNameUgly, "/logging.logwriter.write") &&
			!strings.HasSuffix(methodNameUgly, "/logging.logger.enter") &&
			!strings.HasSuffix(methodNameUgly, "/logging.logger.leave") {
			break
		}
	}
	methodName := strings.Replace(normalize(methodNameUgly), strings.ToLower(name)+".", "", 1)
	if methodName == strings.ToLower(name) || strings.HasPrefix(methodName, "func") {
		methodName = "lambda"
	}

	_msg := ""
	for i := 0; i < len(msg); i++ {
		if _msg != "" {
			_msg += " "
		}
		_msg += fmt.Sprint(msg[i])
	}
	log.Printf("%s [%s#%s] %s\n", level, name, methodName, _msg)
}
