package logger

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type (
	// LogToStream log messages to stream
	LogToStream struct {
		MaxProgress     int
		CurrentProgress int

		LogStream         io.Writer
		ProgressStream    io.Writer
		MaxProgressStream io.Writer
	}
)

// Info send [INFO] message to log
func (log *LogToStream) Info(message string) {
	logMessageToStream(log.LogStream, "Info", message)
}

// Warn send [WARN] message to log
func (log *LogToStream) Warn(message string) {
	logMessageToStream(log.LogStream, "Warn", message)
}

// Error send [ERROR] message to log
func (log *LogToStream) Error(message string) {
	logMessageToStream(log.LogStream, "Error", message)
}

// Fatal send [FATAL] message to log
func (log *LogToStream) Fatal(message string) {
	logMessageToStream(log.LogStream, "Fatal", message)
}

// SetProgressMax sets the maximum Progress in int
func (log *LogToStream) SetProgressMax(max int) {
	var b []byte
	binary.LittleEndian.PutUint64(b, uint64(max))
	log.MaxProgressStream.Write(b)
	log.MaxProgress = max
}

// ProgressInc incease current progress with int as param
func (log *LogToStream) ProgressInc(progress int) {
	var b []byte
	binary.LittleEndian.PutUint64(b, uint64(progress))
	log.ProgressStream.Write(b)
	log.CurrentProgress = log.CurrentProgress + progress
}

// Finalize finalize the execution and max out progress
func (log *LogToStream) Finalize() {
	logMessageToStream(log.LogStream, "Finalize", "Process has finished")
	log.ProgressStream.Write(nil)
}

func logMessageToStream(stream io.Writer, logType, message string) {
	stream.Write([]byte(fmt.Sprintf("[%s] %s", strings.ToUpper(logType), message)))
}
