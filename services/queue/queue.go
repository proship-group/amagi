package queue

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
	"time"

	utils "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/services/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Executor queue execute command
	Executor interface {
		// Execute the method to call during dequeueing
		Execute() error
		// Identity must be an identifying string for the item
		Identity() string
	}

	// Statuses import queue statuses
	Statuses int

	// Queue object of a queue item
	Queue struct {
		ID         bson.ObjectId `bson:"_id"`
		Status     Statuses      `bson:"status"`
		CreatedAt  time.Time     `bson:"created_at"`
		StartedAt  time.Time     `bson:"started_at"`
		FinishedAt time.Time     `bson:"finished_at"`
		ItemData   []byte        `bson:"item_data"`
		ItemType   string        `bson:"item_type"`

		ItemExec                  Executor `bson:"-" json:"-"`
		notFilterQueueNameDequeue bool
	}
)

const (
	// StatusQueued the item is registered to queue
	StatusQueued Statuses = iota
	// StatusProgress the queue item is currently being run
	StatusProgress
	// StatusDone the item has completed execution
	StatusDone
	// StatusError the item encountered an error
	StatusError
)

var (
	// QueueCollection collection name
	QueueCollection = "queue_import"
)

// Enqueue adds the item to queue db
//
// For example:
//
//     go models.Queue{ItemData: d}.Enqueue()
//
func (item Queue) Enqueue() error {
	if item.ItemExec == nil {
		return fmt.Errorf("Queue item must have ItemExec: %v", item)
	}
	var data bytes.Buffer
	en := gob.NewEncoder(&data)
	if err := en.Encode(&item.ItemExec); err != nil {
		utils.Error(fmt.Sprintf("[Amagi-Queue] error Encoding to GOB: %v", err))
		return err
	}
	sc := database.SessionCopy()
	defer sc.Close()
	coll := sc.DB(database.Db).C(QueueCollection)

	item.ItemData = data.Bytes()
	item.ID = bson.NewObjectId()
	item.Status = StatusQueued
	item.CreatedAt = time.Now()
	item.ItemType = item.ExecName()

	if err := coll.Insert(item); err != nil {
		utils.Error(fmt.Sprintf("[Amagi-Queue] error Enqueue: %v", err))
		return err
	}
	return nil
}

// Dequeue finds an item and claims it for processing
//
// For example:
//
//     var queueItem Queue
//     if err := queueItem.Dequeue(); err != nil {}
//     fmt.Print(queueItem.Status)
//
func (item *Queue) Dequeue(typeName string) error {
	sc := database.SessionCopy()
	defer sc.Close()
	coll := sc.DB(database.Db).C(QueueCollection)

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"status":     StatusProgress,
			"started_at": time.Now(),
		}},
		ReturnNew: true,
	}
	// var dequeued *deququed
	selector := bson.M{"status": StatusQueued}
	if !item.notFilterQueueNameDequeue {
		selector["item_type"] = typeName
	}

	if _, err := coll.Find(selector).Sort("created_at").Apply(change, item); err != nil {
		if err != mgo.ErrNotFound {
			// Do not print error message if none found
			utils.Error(fmt.Sprintf("[Amagi-Queue] error Dequeue: %v", err))
		}
		return err
	}
	de := gob.NewDecoder(bytes.NewReader(item.ItemData))
	if err := de.Decode(&item.ItemExec); err != nil {
		utils.Error(fmt.Sprintf("[Amagi-Queue] error Decoding GOB: %v", err))
		return err
	}
	return nil
}

// FilterDequeue set to filter or not queueName during dequeue
// setting this this false may cause `panic`s on runtime, be careful!!!
func (item *Queue) FilterDequeue(filter bool) {
	item.notFilterQueueNameDequeue = !filter
}

// Success sets Queue.Status = StatusDone
func (item *Queue) Success() error {
	if err := item.updateQueue(StatusDone); err != nil {
		utils.Error(fmt.Sprintf("[Amagi-Queue] error ImportSuccess %v", err))
		return err
	}
	return nil
}

// Fail sets Queue.Status = StatusError
func (item *Queue) Fail() error {
	if err := item.updateQueue(StatusError); err != nil {
		utils.Error(fmt.Sprintf("[Amagi-Queue] error ImportFail %v", err))
		return err
	}
	return nil
}

// ExecName get the calculated name of the Executor dataItem
func (item *Queue) ExecName() string {
	return GetTypeName(item.ItemExec)
}

// GetTypeName returns the last index after split by '.'
func GetTypeName(t interface{}) string {
	if t == nil {
		return ""
	}
	l := strings.Split(reflect.TypeOf(t).String(), ".")
	if len(l) <= 0 {
		return ""
	}
	n := l[len(l)-1]
	return n
}

// CleanUp sets the item to nil
func (item *Queue) CleanUp() {
	notFilterQueueNameDequeue := item.notFilterQueueNameDequeue
	item = &Queue{}
	item.notFilterQueueNameDequeue = notFilterQueueNameDequeue
	utils.Info(fmt.Sprintf("[Amagi-Queue] Item cleaned-up: %v", item))
}

func (item *Queue) updateQueue(status Statuses) error {

	sc := database.SessionCopy()
	defer sc.Close()
	coll := sc.DB(database.Db).C(QueueCollection)

	query := bson.M{"_id": item.ID}
	update := bson.M{"$set": bson.M{
		"status":      status,
		"finished_at": time.Now(),
	}}
	if err := coll.Update(query, update); err != nil {
		return err
	}
	return nil
}

func (status Statuses) String() string {
	return []string{
		"StatusQueued",
		"StatusProgress",
		"StatusDone",
	}[status]
}
