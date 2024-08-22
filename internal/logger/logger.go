package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	Info  *log.Logger
	Trace *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func init() {
	infoPrefix := color.BlueString("[INFO] ")
	Info = log.New(os.Stdout, infoPrefix, log.LstdFlags)

	tracePrefix := color.GreenString("[TRACE] ")
	Trace = log.New(os.Stdout, tracePrefix, log.LstdFlags)

	warnPrefix := color.YellowString("[WARNING] ")
	Warn = log.New(os.Stdout, warnPrefix, log.LstdFlags)

	errorPrefix := color.RedString("[ERROR] ")
	Error = log.New(os.Stderr, errorPrefix, log.LstdFlags|log.Lshortfile)
}
