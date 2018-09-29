package amagi

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/b-eee/amagi/api/slack"
	"github.com/b-eee/amagi/helpers"
	"github.com/b-eee/amagi/services/sentry"
	"github.com/k0kubun/pp"
	"gopkg.in/mgo.v2/bson"
)

var (
	logInterface = map[string]string{
		"i": "INFO",
		"w": "WARN",
		"e": "ERROR",
		"f": "FATAL",
	}
)

type (
	// UtilsLogger alias for root functions
	UtilsLogger struct{}
)

// Init initialize slack API
func Init(host slack.Host) {
	slack.Init(host)
}

// Info print to stdout our message
func Info(msg string) error {
	str := fmt.Sprintf("%s %s", timeLoglevel("i"), FgColorizer(msg, "default"))
	fmt.Println(str)

	return nil
}

// Warn print to stdout
func Warn(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("w"), FgColorizer(msg, "default"))
	fmt.Println(str)
}

// Error print to stdout
func Error(msg string) error {
	str := fmt.Sprintf("%s %s", timeLoglevel("e"), FgColorizer(msg, "default"))
	fmt.Println(str)

	sentry.SendToSentry(msg)
	return fmt.Errorf(str)
}

// Fatal fatal print to stdout
func Fatal(msg string) {
	str := errMsgFmt("f", FgColorizer(msg, "default"))

	go slack.Send("", str)

	sentry.SendToSentry(msg)
	fmt.Println(str)
}

// Pretty Printer for DEBUG
func Pretty(obj interface{}, msg string) {
	str := fmt.Sprintf("--- %s ---", msg)

	fmt.Println(str)
	pp.Println(obj)
}

// PrettyBson Pretty Printer for []bson
func PrettyBson(slice []bson.M, msg string) {
	helpers.PrintBsonSlice(slice, msg)
}

// ExceptionDump start watching stack trace
func ExceptionDump() {
	if e := recover(); e != nil {
		DumpStack(e, debug.Stack())
		panic(e)
	}
}

// DumpStack dump stack trace
func DumpStack(e interface{}, stack []byte) {

	if err := slack.Send(e, string(stack)); err != nil {
		fmt.Printf("cant send to slack %v\n", err)
	}

	sentry.SendToSentry(string(stack))
	fmt.Printf("[%v] crashing...\n", PrettyPrintTime(time.Now()))
}

// UTILS
// LogLevel construct log level msg
func logLevel(key string) string {
	str := fmt.Sprintf("%s", FgColorizer("["+logInterface[key]+"]", key))
	return str
}

func timeLoglevel(logLevelStr string) string {
	str := fmt.Sprintf("%s %s", FgColorizer("["+time.Now().Format(time.RFC3339Nano)+"]", "default"), logLevel(logLevelStr))
	return str
}

// PrettyPrintTime pretty print a time value to readable
func PrettyPrintTime(timeVal time.Time) string {
	return timeVal.Format("Mon Jan _2 15:04:05 2006")
}

func errMsgFmt(logLevel, msg string) string {
	return fmt.Sprintf("%s %s", timeLoglevel(logLevel), msg)
}

// Initialize initialize the logger with the ID
func (log *UtilsLogger) Initialize(id string) {}

// Info send [INFO] message to log
func (log *UtilsLogger) Info(message string) {
	Info(message)
}

// Warn send [WARN] message to log
func (log *UtilsLogger) Warn(message string) {
	Warn(message)
}

// Error send [ERROR] message to log
func (log *UtilsLogger) Error(message string) {
	Error(message)
}

// Fatal send [FATAL] message to log
func (log *UtilsLogger) Fatal(message string) {
	Fatal(message)
}

// SetProgressMax sets the maximum Progress in int
func (log *UtilsLogger) SetProgressMax(max int) {}

// ProgressInc incease current progress with int as param
func (log *UtilsLogger) ProgressInc(progress int) {}

// Finalize finalize the execution and max out progress
func (log *UtilsLogger) Finalize() {}
