package logger

import (
	"context"
	"fmt"
	"strings"
	"time"

	mongodb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// LogToMongo log to mongodb logger
	LogToMongo struct {
		ID              bson.ObjectId `bson:"_id"`
		MaxProgress     int           `bson:"max_progress"`
		CurrentProgress int           `bson:"progress" json:"progress"`

		Database       string `bson:"-" json:"-"`
		CollectionName string `bson:"-" json:"-"`

		Session    *mongodb.Session    `bson:"-" json:"-"`
		Collection *mongodb.Collection `bson:"-" json:"-"`

		PersistProgressIntervalMS int `bson:"-" json:"-"`

		persistedProgress int
		persistentCloser  context.CancelFunc
	}

	// LogMessage log messages
	LogMessage struct {
		ID        bson.ObjectId
		Message   string
		CreatedAt time.Time
		LogID     string
	}
)

var (
	// PersistProgressIntervalMSDefault default interval value
	PersistProgressIntervalMSDefault = 2000
)

// InitializeCollection initialize collection
func (log *LogToMongo) InitializeCollection() {
	if log.Collection != nil {
		return
	}

	log.ID = bson.NewObjectId()
	log.Collection = log.Session.DB(log.Database).C(log.CollectionName)
	if log.PersistProgressIntervalMS <= 0 {
		log.PersistProgressIntervalMS = PersistProgressIntervalMSDefault
	}
	err := log.Collection.Insert(log)
	if err != nil {
		fmt.Println(fmt.Sprintf("error LogToMongo.InitializeCollection insert: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	log.persistentCloser = cancel
	go log.persistProgress(ctx)
}

// Info send [INFO] message to log
func (log *LogToMongo) Info(message string) {
	log.InitializeCollection()
	createLogMessage(log.Collection, log.ID.Hex(), "Info", message)
}

// Warn send [WARN] message to log
func (log *LogToMongo) Warn(message string) {
	log.InitializeCollection()
	createLogMessage(log.Collection, log.ID.Hex(), "Warn", message)
}

// Error send [ERROR] message to log
func (log *LogToMongo) Error(message string) {
	log.InitializeCollection()
	createLogMessage(log.Collection, log.ID.Hex(), "Error", message)
}

// Fatal send [FATAL] message to log
func (log *LogToMongo) Fatal(message string) {
	log.InitializeCollection()
	createLogMessage(log.Collection, log.ID.Hex(), "Fatal", message)
}

// SetProgressMax sets the maximum Progress in int
func (log *LogToMongo) SetProgressMax(max int) {
	log.InitializeCollection()
	log.MaxProgress = max

	err := log.Collection.Update(
		bson.M{"_id": log.ID},
		bson.M{"$set": bson.M{"max_progress": log.MaxProgress}},
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("error LogToMongo.SetProgressMax: %v", err))
	}
}

// ProgressInc incease current progress with int as param
func (log *LogToMongo) ProgressInc(progress int) {
	log.InitializeCollection()
	log.CurrentProgress = log.CurrentProgress + progress
}

// Finalize finalize the execution and max out progress
func (log *LogToMongo) Finalize() {
	defer log.Session.Close()
	log.persistentCloser()
	createLogMessage(log.Collection, log.ID.Hex(), "Finalize", "Process has finished")
}

func createLogMessage(collection *mongodb.Collection, logID, logType, message string) {
	err := collection.Insert(
		LogMessage{
			ID:        bson.NewObjectId(),
			Message:   fmt.Sprintf("[%s] %s", strings.ToUpper(logType), message),
			CreatedAt: time.Now(),
			LogID:     logID,
		},
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("error LogToMongo.%s: %v", logType, err))
	}
}

func (log *LogToMongo) persistProgress(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := log.Collection.Update(
				bson.M{"_id": log.ID},
				bson.M{"$set": bson.M{"progress": log.CurrentProgress}},
			)
			if err != nil {
				fmt.Println(fmt.Sprintf("error LogToMongo.SetProgressMax: %v", err))
			}
		}
		time.Sleep(time.Duration(log.PersistProgressIntervalMS) * time.Millisecond)
	}
}
