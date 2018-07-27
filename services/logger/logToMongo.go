package logger

import (
	mongodb "gopkg.in/mgo.v2"
)

type (
	// LogToMongo log to mongodb logger
	LogToMongo struct {
		Session         *mongodb.Session
		MaxProgress     int
		CurrentProgress int
		Identity        string

		Database   string
		Collection string
	}
)

// SetIdentity set Identity
func (log *LogToMongo) SetIdentity(ident string) {
	log.Identity = ident
}

// Info send [INFO] message to log
func (log *LogToMongo) Info(string) {
	col := log.Session.DB(log.Database).C(log.Collection)
	_ = col
}

// Warn send [WARN] message to log
func (log *LogToMongo) Warn(string) {

}

// Error send [ERROR] message to log
func (log *LogToMongo) Error(string) {

}

// Fatal send [FATAL] message to log
func (log *LogToMongo) Fatal(string) {

}

// SetProgressMax sets the maximum Progress in int
func (log *LogToMongo) SetProgressMax(int) {

}

// ProgressInc incease current progress with int as param
func (log *LogToMongo) ProgressInc(int) {

}

// Finalize finalize the execution and max out progress
func (log *LogToMongo) Finalize() {

}
