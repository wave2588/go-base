package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	eParser "github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	logNameDebug   = "debug.log"
	logNameInfo    = "info.log"
	logNameWarning = "warn.log"
	logNameError   = "error.log"
	logNameRequest = "request.log"

	logTimeFormat = "2006-01-02 15:04:05.000"
)

// Logger global logger.
type Logger struct {
	Debug *logrus.Logger
	Info  *logrus.Logger
	Warn  *logrus.Logger
	Error *logrus.Logger

	Request *lumberjack.Logger
}

// RequestLog contains request log info.
type RequestLog struct {
	During time.Duration `json:"during"`
	Method string        `json:"method"`
	Status int           `json:"status"`
	IP     string        `json:"ip"`
	URI    string        `json:"uri"`
	Msg    string        `json:"msg"`
}

var (
	logger    Logger
	logPath   = "./logs"
	debug     bool
	logPrefix string

	filePathPrefix string
)

// Init creates global logger instances.
func Init(debugMode bool) {
	err := os.MkdirAll(logPath, 0o700)
	if err != nil {
		panic(err)
	}

	debug = debugMode

	filePathPrefix, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	logger = Logger{
		Debug:   newLogger(logNameDebug, logrus.DebugLevel),
		Info:    newLogger(logNameInfo, logrus.InfoLevel),
		Warn:    newLogger(logNameWarning, logrus.WarnLevel),
		Error:   newLogger(logNameError, logrus.ErrorLevel),
		Request: newLogWriter(path.Join(logPath, logNameRequest)),
	}
}

// SetPrefix sets the output prefix for the logger.
func SetPrefix(prefix string) {
	logPrefix = prefix
}

func newLogger(fileName string, level logrus.Level) *logrus.Logger {
	fileName = path.Join(logPath, fileName)

	l := &logrus.Logger{
		Out:       nil,
		Formatter: new(logFormatter),
		Level:     level,
		Hooks:     nil,
	}

	if debug {
		l.SetOutput(io.MultiWriter(os.Stdout, newLogWriter(fileName)))
		return l
	}

	if level >= logrus.DebugLevel {
		l.SetOutput(io.Discard)
	} else {
		l.SetOutput(io.MultiWriter(os.Stdout, newLogWriter(fileName)))
	}

	return l
}

func newLogWriter(logPath string) *lumberjack.Logger {
	logger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    30,
		MaxBackups: 100,
		MaxAge:     30,
	}

	return logger
}

func DiscardBuiltinLogOutput() {
	log.SetOutput(io.Discard)
}

// GetErrorLogger returns error logger
func GetErrorLogger() *logrus.Logger {
	return logger.Error
}

// logFormatter defines custom formatter for logrus
type logFormatter struct{}

// Format formats log output.
func (f *logFormatter) Format(e *logrus.Entry) ([]byte, error) {
	format := ""

	if logPrefix != "" {
		format = fmt.Sprintf("%s [%s][%s] %s", e.Time.Format(logTimeFormat), logPrefix, e.Level.String(), e.Message)
		return []byte(format), nil
	}

	format = fmt.Sprintf("%s [%s] %s", e.Time.Format(logTimeFormat), e.Level.String(), e.Message)
	return []byte(format), nil
}

// Debugf logs in Debug level.
func Debugf(format string, v ...interface{}) {
	logger.Debug.Debug(logHandler(format, v))
}

// Debug logs in Debug level.
func Debug(v ...interface{}) {
	logger.Debug.Debug(logHandler("", v))
}

// DebugSQL logs SQL statement in Debug level.
func DebugSQL(sql string, omitList *[]string) {
	_, file, line, _ := runtime.Caller(2)
	filePath := fmt.Sprintf("%s:%d", file, line)

	if omitList != nil {
		for _, keyword := range *omitList {
			if strings.HasSuffix(filePath, keyword) {
				return
			}
		}
	}

	msg := fmt.Sprintf("[%s] %s", fileInfo(), sql)
	logger.Debug.Debug(msg)
}

// Infof logs in Info level.
func Infof(format string, v ...interface{}) {
	logger.Info.Info(logHandler(format, v))
}

// Info logs in Info level.
func Info(v ...interface{}) {
	logger.Info.Info(logHandler("", v))
}

// Warnf logs in Warn level.
func Warnf(format string, v ...interface{}) {
	logger.Warn.Warn(logHandler(format, v))
}

// Warn logs in Warn level.
func Warn(v ...interface{}) {
	logger.Warn.Warn(logHandler("", v))
}

// Errorf logs in Error level.
func Errorf(format string, v ...interface{}) {
	logger.Error.Error(logHandler(format, v))
}

// Error logs in Error level.
func Error(v ...interface{}) {
	logger.Error.Error(logHandler("", v))
}

func ErrorSQL(query interface{}, args []interface{}, err error) {
	var rawQ string
	if _q, ok := query.([]string); ok {
		rawQ = strings.Join(_q, "\n")
	} else {
		rawQ = fmt.Sprintf("%s", query)
	}

	argsIdx := 0
	var q strings.Builder
	for _, ch := range rawQ {
		if ch == '?' {
			if argsIdx < len(args) {
				q.WriteString(fmt.Sprintf("%v", args[argsIdx]))
				argsIdx++
			}
		} else {
			q.WriteRune(ch)
		}
	}

	logger.Error.Error(logHandler("\n[query]\n%s\n[error]\n%v", []interface{}{q.String(), err}))
}

// Fatalf logs in Fatal level.
func Fatalf(format string, v ...interface{}) {
	msg := logHandler(format, v)
	logger.Error.Fatal(msg)
}

// Fatal logs in Fatal level.
func Fatal(v ...interface{}) {
	msg := logHandler("", v)
	logger.Error.Fatal(msg)
}

// Panicf logs in Panic level.
func Panicf(format string, v ...interface{}) {
	msg := logHandler(format, v)
	logger.Error.Panicf(msg)
}

// Panic logs in Panic level.
func Panic(v ...interface{}) {
	msg := logHandler("", v)
	logger.Error.Panic(msg)
}

// LogRequest logs api request.
func LogRequest(reqLog RequestLog) {
	v := []interface{}{}
	v = append(v, time.Now().Format(logTimeFormat))
	v = append(v, fmt.Sprintf("[%s]", reqLog.Method))
	v = append(v, reqLog.Status)
	v = append(v, reqLog.During)
	v = append(v, reqLog.IP)
	v = append(v, reqLog.URI)
	v = append(v, reqLog.Msg)

	//      time [method] [status] during ip uri msg
	format := "%v %-9s [%d] %-13v %-16v %s %s\n"

	msg := fmt.Sprintf(format, v...)
	logger.Request.Write([]byte(msg))
}

func logHandler(format string, v []interface{}) (msg string) {
	defer func() {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
	}()

	if debug {
		msg = fmt.Sprintf("[%s] ", fileInfo())
	}

	if v == nil {
		return msg + format
	}

	for i := 0; i < len(v); i++ {
		v[i] = extract(v[i])
	}

	if format == "" {
		for _, v := range v {
			msg += fmt.Sprint(v)
		}
		return msg
	}

	return msg + fmt.Sprintf(format, v...)
}

func extract(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	if e, ok := v.(error); ok {
		err := eParser.Wrap(e, 3)
		return fmt.Sprintf("%s\n%s", err.Error(), string(err.Stack()))
	}

	t := reflect.TypeOf(v)

	if strings.HasPrefix(t.String(), "*big") ||
		strings.HasPrefix(t.String(), "enums") {
		return v
	}

	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}

	switch t.Kind() {
	case reflect.Struct:
		b, err := json.Marshal(v)
		if err != nil {
			e := eParser.Wrap(err, 0)
			return fmt.Sprintf("%s\n%s", e.Error(), string(e.Stack()))
		}
		return string(b)
	case reflect.Ptr:
		if reflect.ValueOf(v).IsNil() {
			return extract(nil)
		}

		return extract(reflect.ValueOf(v).Elem().Interface())
	default:
		return v
	}
}

func fileInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "<???>"
	} else {
		var i int
		for i = 0; i < len(file) && i < len(filePathPrefix); i++ {
			if file[i] != filePathPrefix[i] {
				break
			}
		}
		commonPrefix := file[:i]
		file = strings.TrimPrefix(file, commonPrefix)
		file = strings.TrimPrefix(file, "/")
	}

	return fmt.Sprintf("%s:%d", file, line)
}
