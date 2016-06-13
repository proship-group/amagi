package hokutoseiUtils

import (
	"fmt"
	"time"
)

var (
	logInterface = map[string]string{
		"i": "INFO",
		"e": "ERROR",
	}
)

// Info print to stdout our message
func Info(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("i"), msg)
	fmt.Println(str)
}

// Error print to stdout
func Error(msg string) {
	str := fmt.Sprintf("%s %s", timeLoglevel("e"), msg)
	fmt.Println(str)
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
