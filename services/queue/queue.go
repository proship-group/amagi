package queue

import (
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
		ItemData   Executor      `bson:"item_data"`
		ItemType   string        `bson:"item_type"`

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

// SetBSON implements bson.Setter. This is needed because mgo does cannot map bson values to interfaces with methods
func (item *Queue) SetBSON(raw bson.Raw) error {

	decoded := new(struct {
		ID         bson.ObjectId `bson:"_id"`
		Status     Statuses      `bson:"status"`
		CreatedAt  time.Time     `bson:"created_at"`
		StartedAt  time.Time     `bson:"started_at"`
		FinishedAt time.Time     `bson:"finished_at"`
		ItemData   interface{}   `bson:"item_data"`
		ItemType   string        `bson:"item_type"`
	})

	bsonErr := raw.Unmarshal(decoded)
	if bsonErr == nil {
		// need to manually assign values here...
		item.ID = decoded.ID
		item.Status = decoded.Status
		item.CreatedAt = decoded.CreatedAt
		item.StartedAt = decoded.StartedAt
		item.FinishedAt = decoded.FinishedAt
		item.ItemType = decoded.ItemType
	} else {
		return bsonErr
	}
	// after setting everything else, get the interface data and assign it to the initialized one in `item`
	// marshal the bson data, then unmarshal it again into the `item.ItemData`
	bsonBytes, _ := bson.Marshal(decoded.ItemData)
	if bsonErr = bson.Unmarshal(bsonBytes, item.ItemData); bsonErr != nil {
		return bsonErr
	}
	return nil
}

// Enqueue adds the item to queue db
//
// For example:
//
//     go models.Queue{ItemData: d}.Enqueue()
//
func (item Queue) Enqueue() error {
	if item.ItemData == nil {
		return fmt.Errorf("Queue item must have ItemData: %v", item)
	}
	sc := database.SessionCopy()
	defer sc.Close()
	coll := sc.DB(database.Db).C(QueueCollection)

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
func (item *Queue) Dequeue() error {
	if item.ItemData == nil {
		return fmt.Errorf("Queue item must have ItemData not set to nil")
	}
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
		selector["item_type"] = item.ExecName()
	}

	if _, err := coll.Find(selector).Sort("created_at").Apply(change, item); err != nil {
		if err != mgo.ErrNotFound {
			// Do not print error message if none found
			utils.Error(fmt.Sprintf("[Amagi-Queue] error Dequeue: %v", err))
		}
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
	if item.ItemData == nil {
		return ""
	}
	l := strings.Split(reflect.TypeOf(item.ItemData).String(), ".")
	if len(l) <= 0 {
		return ""
	}
	n := l[len(l)-1]
	return n
}

// CleanUp sets the item to nil
func (item *Queue) CleanUp() {
	item.ID = ""
	item.Status = 0
	item.CreatedAt = time.Time{}
	item.StartedAt = time.Time{}
	item.FinishedAt = time.Time{}
	item.ItemType = ""
	p := reflect.ValueOf(item.ItemData).Elem()
	p.Set(reflect.Zero(p.Type()))
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
