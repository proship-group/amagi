package logger

import (
	"sync"
)

type (
	// Logger copied from queue to eliminate circular import
	Logger interface {
		Initialize(string)
		Info(string)
		Warn(string)
		Error(string)
		Fatal(string)
		SetProgressMax(int)
		ProgressInc(int)
		Finalize()
	}

	// LogToLoggers log to list of loggers
	LogToLoggers struct {
		Loggers []Logger
	}
)

// Initialize initialize the logger with the ID
func (log *LogToLoggers) Initialize(id string) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.Initialize(id)
		}(logger)
	}
	wg.Wait()
}

// Info send [INFO] message to log
func (log *LogToLoggers) Info(message string) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.Info(message)
		}(logger)
	}
	wg.Wait()
}

// Warn send [WARN] message to log
func (log *LogToLoggers) Warn(message string) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.Warn(message)
		}(logger)
	}
	wg.Wait()
}

// Error send [ERROR] message to log
func (log *LogToLoggers) Error(message string) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.Error(message)
		}(logger)
	}
	wg.Wait()
}

// Fatal send [FATAL] message to log
func (log *LogToLoggers) Fatal(message string) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.Fatal(message)
		}(logger)
	}
	wg.Wait()
}

// SetProgressMax sets the maximum Progress in int
func (log *LogToLoggers) SetProgressMax(max int) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.SetProgressMax(max)
		}(logger)
	}
	wg.Wait()
}

// ProgressInc incease current progress with int as param
func (log *LogToLoggers) ProgressInc(progress int) {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.ProgressInc(progress)
		}(logger)
	}
	wg.Wait()
}

// Finalize finalize the execution and max out progress
func (log *LogToLoggers) Finalize() {
	var wg *sync.WaitGroup
	wg.Add(len(log.Loggers))
	for _, logger := range log.Loggers {
		go func(logger Logger) {
			defer wg.Done()
			logger.Finalize()
		}(logger)
	}
	wg.Wait()
}
