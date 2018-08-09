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
		ParentID        bson.ObjectId `bson:"task_id" json:"task_id"`

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
		ID        bson.ObjectId `bson:"_id"`
		Message   string        `bson:"message"`
		CreatedAt time.Time     `bson:"created_at"`
		ParentID  bson.ObjectId `bson:"parent_id"`
	}
)

var (
	// PersistProgressIntervalMSDefault default interval value
	PersistProgressIntervalMSDefault = 2000
	// SoftTimeout persistence go routine soft timeout
	SoftTimeout = 20 * time.Minute
	// HardTimeout persistence go routine hard timeout
	HardTimeout = 25 * time.Minute
)

// Initialize initialize collection
func (log *LogToMongo) Initialize(id string) {
	if log.Collection != nil {
		return
	}
	if log.persistentCloser != nil {
		log.persistentCloser()
	}
	log.ID = bson.NewObjectId()
	log.ParentID = bson.ObjectIdHex(id)
	log.Collection = log.Session.DB(log.Database).C(log.CollectionName)
	if log.PersistProgressIntervalMS <= 0 {
		log.PersistProgressIntervalMS = PersistProgressIntervalMSDefault
	}
	err := log.Collection.Insert(log)
	if err != nil {
		fmt.Println(fmt.Sprintf("error LogToMongo.InitializeCollection insert: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), SoftTimeout)
	log.persistentCloser = cancel
	go log.persistProgress(ctx, HardTimeout)
}

// Info send [INFO] message to log
func (log *LogToMongo) Info(message string) {
	createLogMessage(log.Collection, log.ID.Hex(), "Info", message)
}

// Warn send [WARN] message to log
func (log *LogToMongo) Warn(message string) {
	createLogMessage(log.Collection, log.ID.Hex(), "Warn", message)
}

// Error send [ERROR] message to log
func (log *LogToMongo) Error(message string) {
	createLogMessage(log.Collection, log.ID.Hex(), "Error", message)
}

// Fatal send [FATAL] message to log
func (log *LogToMongo) Fatal(message string) {
	createLogMessage(log.Collection, log.ID.Hex(), "Fatal", message)
}

// SetProgressMax sets the maximum Progress in int
func (log *LogToMongo) SetProgressMax(max int) {
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
	log.CurrentProgress = log.CurrentProgress + progress
}

// Finalize finalize the execution and max out progress
func (log *LogToMongo) Finalize() {
	// if log.Session != nil {
	// 	defer log.Session.Close()
	// }
	if log.persistentCloser != nil {
		log.persistentCloser()
	}
}

func createLogMessage(collection *mongodb.Collection, logID, logType, message string) {
	err := collection.Insert(
		LogMessage{
			ID:        bson.NewObjectId(),
			Message:   fmt.Sprintf("[%s] %s", strings.ToUpper(logType), message),
			CreatedAt: time.Now(),
			ParentID:  bson.ObjectIdHex(logID),
		},
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("error LogToMongo.%s: %v", logType, err))
	}
}

func (log *LogToMongo) persistProgress(ctx context.Context, hardTimeout time.Duration) {
	update := func(log *LogToMongo) {
		err := log.Collection.Update(
			bson.M{"_id": log.ID},
			bson.M{"$set": bson.M{"progress": log.CurrentProgress}},
		)
		if err != nil {
			fmt.Println(fmt.Sprintf("error LogToMongo.SetProgressMax: %v", err))
		}
	}

	for {
		select {
		case <-time.After(hardTimeout):
			update(log)
			log = nil
			return
		case <-ctx.Done():
			update(log)
			log = nil
			return
		default:
			update(log)
		}
		time.Sleep(time.Duration(log.PersistProgressIntervalMS) * time.Millisecond)
	}
}
