package amagi

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/b-eee/amagi/api/pubnub"
	"github.com/b-eee/amagi/api/slack"
)

var (
	logInterface = map[string]string{
		"i": "INFO",
		"e": "ERROR",
		"f": "FATAL",
	}
)

// Init initialize slack API
func Init(host slack.Host) {
	fmt.Println("initializing slack...")
	slack.Init(host)

	pubnub.SetPubNubConnection()
}

// Info print to stdout our message
func Info(msg string) {
	var wg sync.WaitGroup
	wg.Add(1)
	str := fmt.Sprintf("%s %s", timeLoglevel("i"), msg)
	fmt.Println(str)

	go func() {
		channel := pubnub.ChanName([]string{"log", "stream"}...)
		message := formatHostName(str, slack.GetMicroAppName(), slack.GetCurrentConfiguredHost())

		// pubnub.Publish(channel, message, &wg)
		pubnub.Publish(channel, message, &wg)
	}()
	wg.Wait()
}

// Error print to stdout
func Error(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("e"), msg)

	fmt.Println(str)
}

// Fatal fatal print to stdout
func Fatal(msg string) {
	str := errMsgFmt("f", msg)

	go slack.Send("", str)
	fmt.Println(str)
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
