package logging

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"runtime"
	"strings"
)

var logCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "mailassistant",
	Subsystem: "logging",
	Name:      "log",
	Help:      "counter per logger and level",
}, []string{"name", "level"})

// NewLogger is a factory for a new log instance with an autodetected log name
func NewLogger() Logger {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	return NewNamedLogger(details.Name())
}

// NewNamedLogger is a factory for a new log instance with an given name
func NewNamedLogger(name string) Logger {
	return newLogger(normalize(name))
}

// NewGlobalLogger is a factory for a new log instance without name and parent
func NewGlobalLogger() Logger {
	return NewNamedLogger("global")
}

// NewNamedLogWriter is a factory for a new log instance with given name as io.Writer
func NewNamedLogWriter(name string) io.Writer {
	return logWriter{NewNamedLogger(name)}
}

type logWriter struct {
	Logger
}

// Writer is converting the byte array to string and delegate to the Debug method
func (l logWriter) Write(p []byte) (n int, err error) {
	msg := "\n" + string(p)
	l.Debug(strings.TrimSuffix(msg, "\n"))
	return len(p), nil
}

var loggerRegistry map[string]logLevel

func init() {
	prometheus.MustRegister(logCounter)
	SetLevel("*", "")
}

// SetLevel is a key, value pair to set a logLevel programmatically
func SetLevel(name, level string) {
	if name == "*" {
		loggerRegistry = make(map[string]logLevel, 0)
	} else {
		loggerRegistry[name] = stringToLogLevel(level)
	}
}

func newLogger(name string) Logger {
	return &iLogger{name}
}

// Logger is a interface witch represents a log implementation
type Logger interface {
	Name() string
	Debug(msg ...interface{})
	Debugf(format string, msg ...interface{})
	Info(msg ...interface{})
	Infof(format string, msg ...interface{})
	Warn(msg ...interface{})
	Warnf(format string, msg ...interface{})
	Severe(msg ...interface{})
	Severef(format string, msg ...interface{})
	Panic(msg ...interface{})
	Enter()
	Leave()
}

// Logger represents a Log Entity
type iLogger struct {
	name string
}

// Name returns the name of the current entity
func (l iLogger) Name() string {
	return l.name
}

// Severe is writing a severe message
func (l iLogger) Severe(msg ...interface{}) {
	if l.isLogLevel(severe) {
		logger(l.name, "SEVERE ", msg)
	}
}

// Severef is writing a severe message
func (l iLogger) Severef(format string, msg ...interface{}) {
	if l.isLogLevel(severe) {
		l.Severe(fmt.Sprintf(format, msg...))
	}
}

// Warn is writing a WARNING message
func (l iLogger) Warn(msg ...interface{}) {
	if l.isLogLevel(warn) {
		logger(l.name, "WARNING", msg)
	}
}

// Warnf is writing a WARNING message
func (l iLogger) Warnf(format string, msg ...interface{}) {
	if l.isLogLevel(warn) {
		l.Warn(fmt.Sprintf(format, msg...))
	}
}

// Info is writing a INFO message
func (l iLogger) Info(msg ...interface{}) {
	if l.isLogLevel(info) {
		logger(l.name, "INFO   ", msg)
	}
}

// Infof is writing a INFO message
func (l iLogger) Infof(format string, msg ...interface{}) {
	if l.isLogLevel(info) {
		l.Info(fmt.Sprintf(format, msg...))
	}
}

// Debug is writing a debug message
func (l iLogger) Debug(msg ...interface{}) {
	if l.isLogLevel(debug) {
		logger(l.name, "DEBUG  ", msg)
	}
}

// Debugf is writing a debug message
func (l iLogger) Debugf(format string, msg ...interface{}) {
	if l.isLogLevel(debug) {
		l.Debug(fmt.Sprintf(format, msg...))
	}
}

// Enter is writing a special debug message that represents to enter a method
func (l iLogger) Enter() {
	l.Debug(">>")
}

// Leave is writing a special debug message that represents to leave a method
func (l iLogger) Leave() {
	l.Debug("<<")
}

// IsLogLevel check the loggerRegistry for a given log level
func (l iLogger) isLogLevel(level logLevel) bool {
	name := strings.Split(l.name, "%")[0]
	if lvl, ok := loggerRegistry[name]; ok && lvl != notExists {
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
func (l iLogger) Panic(msg ...interface{}) {
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
		},
	}

	makeItNicer := name
	for _, modify := range modifyLogic {
		makeItNicer = modify(makeItNicer)
	}

	return makeItNicer
}

func logger(name, level string, msg []interface{}) {
	logCounter.WithLabelValues(name, strings.TrimSpace(level)).Inc()
	methodNameUgly := ""
	for i := 2; i < 5; i++ {
		pc, _, _, _ := runtime.Caller(i)
		method := runtime.FuncForPC(pc)
		methodNameUgly = strings.ToLower(method.Name())
		if !strings.Contains(methodNameUgly, "/logging.ilogger") &&
			!strings.Contains(methodNameUgly, "logging.ilogger") &&
			!strings.Contains(methodNameUgly, "logging.ilogger.panic") &&
			!strings.HasSuffix(methodNameUgly, "logging.logwriter.write") {
			break
		}
	}
	methodName := strings.Replace(normalize(methodNameUgly), strings.ToLower(name)+".", "", 1)
	if methodName == strings.ToLower(name) ||
		strings.ToLower(methodName) == strings.ToLower(name) ||
		strings.HasPrefix(methodName, "func") {
		methodName = "lambda"
	}
	for _, fcc := range []func(string) string{
		func(ugly string) string { return strings.TrimPrefix(ugly, "mailassistant.actions.") },
		func(ugly string) string { return strings.TrimPrefix(ugly, "mailAssistant.actions.") }} {
		methodName = fcc(methodName)
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
