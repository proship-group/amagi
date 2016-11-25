package amagi

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/b-eee/amagi/api/pubnub"
	"github.com/b-eee/amagi/api/slack"
	"github.com/b-eee/amagi/services/sentry"

	"github.com/k0kubun/pp"
)

var (
	logInterface = map[string]string{
		"i": "INFO",
		"w": "WARN",
		"e": "ERROR",
		"f": "FATAL",
	}
)

// Init initialize slack API
func Init(host slack.Host) {
	slack.Init(host)

	pubnub.SetPubNubConnection()
}

// Info print to stdout our message
func Info(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("i"), msg)
	fmt.Println(str)

	var wg sync.WaitGroup
	wg.Add(1)
	go pubnub.Publish(str, &wg)
}

// Warn print to stdout
func Warn(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("w"), msg)

	fmt.Println(str)
	var wg sync.WaitGroup
	wg.Add(1)
	go pubnub.Publish(str, &wg)
}

// Error print to stdout
func Error(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("e"), msg)

	fmt.Println(str)
	sentry.SendToSentry(msg)
	var wg sync.WaitGroup
	wg.Add(1)
	go pubnub.Publish(str, &wg)
}

// Fatal fatal print to stdout
func Fatal(msg string) {
	str := errMsgFmt("f", msg)

	go slack.Send("", str)
	var wg sync.WaitGroup
	wg.Add(1)
	go pubnub.Publish(str, &wg)

	sentry.SendToSentry(msg)
	fmt.Println(str)
}

// Pretty Printer for DEBUG
func Pretty(obj interface{}, msg string) {
	str := fmt.Sprintf("--- %s ---", msg)

	fmt.Println(str)
	pp.Println(obj)
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
	str := fmt.Sprintf("[%s]", logInterface[key])
	return str
}

func timeLoglevel(logLevelStr string) string {
	str := fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC822Z), logLevel(logLevelStr))
	return str
}

// PrettyPrintTime pretty print a time value to readable
func PrettyPrintTime(timeVal time.Time) string {
	return timeVal.Format("Mon Jan _2 15:04:05 2006")
}

func errMsgFmt(logLevel, msg string) string {
	return fmt.Sprintf("%s %s", timeLoglevel(logLevel), msg)
}
