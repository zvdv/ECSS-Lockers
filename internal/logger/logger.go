package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	info  *log.Logger
	trace *log.Logger
	warn  *log.Logger
	error *log.Logger
)

func init() {
	infoPrefix := color.BlueString("[INFO] ")
	info = log.New(os.Stdout, infoPrefix, log.LstdFlags)

	tracePrefix := color.GreenString("[TRACE] ")
	trace = log.New(os.Stdout, tracePrefix, log.LstdFlags)

	warnPrefix := color.YellowString("[WARNING] ")
	warn = log.New(os.Stdout, warnPrefix, log.LstdFlags)

	errorPrefix := color.RedString("[ERROR] ")
	error = log.New(os.Stderr, errorPrefix, log.LstdFlags|log.Llongfile)
}

func Info(format string, v ...any) {
	info.Printf(format+"\n", v...)
}

func Trace(format string, v ...any) {
	trace.Printf(format+"\n", v...)
}

func Warn(format string, v ...any) {
	warn.Printf(format+"\n", v...)
}

func Error(format string, v ...any) {
	error.Printf(format+"\n", v...)
}

func Fatal(v ...any) {
	error.Fatal(v...)
}
