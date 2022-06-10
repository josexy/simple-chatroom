package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Level int

const (
	OFF Level = iota
	INFO
	WARN
	ERROR
	FATAL
	DEBUG
)

var DebugMode = true
var StdOutput = color.Output

func typeOf(level Level) (typ string) {
	switch level {
	case INFO:
		typ = Green("[INFO]")
	case WARN:
		typ = Yellow("[WARN]")
	case ERROR:
		typ = Red("[ERROR]")
	case FATAL:
		typ = HiRed("[FATAL]")
	case DEBUG:
		typ = HiCyan("[DEBUG]")
	}
	return
}

// printf example: [INFO] [main.go:10#funcname] this is a log message
func printf(level Level, format string, v ...interface{}) {
	if level == DEBUG && !DebugMode {
		return
	}
	pc, file, line, _ := runtime.Caller(2)
	ls := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	funcName := ls[len(ls)-1]

	_, err := fmt.Fprintf(StdOutput, "%s [%s:%s#%s] %s\n",
		typeOf(level),
		Green(filepath.Base(file)), Red(strconv.Itoa(line)), Yellow(funcName),
		fmt.Sprintf(format, v...))
	if err != nil {
		return
	}
	if level == FATAL {
		os.Exit(0)
	}
}

func Debug(format string, v ...interface{}) {
	printf(DEBUG, format, v...) // runtime.Caller(1)
}

func Info(format string, v ...interface{}) {
	printf(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	printf(WARN, format, v...)
}

func ErrorBy(err error) {
	Error("%s", err.Error())
}

func Error(format string, v ...interface{}) {
	printf(ERROR, format, v...)
}

func Fatal(format string, v ...interface{}) {
	printf(FATAL, format, v...)
}
