package scheduler

import (
	"labix.org/v2/mgo/bson"
	"sync"
	"time"
)

type Job interface {
	GetId() bson.ObjectId
	Run()
	Deadline() time.Time
}

type job struct {
	Data     Job   // The value of the item; any struct that implement Job
	priority int64 // The priority of the item in the queue.(date to be executed in unixtimestamp)
	index    int   // The index of the item in the heap.
	sync.Mutex
}
