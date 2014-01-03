package scheduler

import (
	"container/heap"
	"labix.org/v2/mgo/bson"
	"log"
	"sync"
	"time"
)

type IJob interface {
	GetId() bson.ObjectId
	Run()
	Deadline() time.Time
}

type job struct {
	Data     IJob  // The value of the item; any struct that implement IJob
	priority int64 // The priority of the item in the queue.
	index    int   // The index of the item in the heap.
	sync.Mutex
}

type delayedJobs []*job

var mu sync.RWMutex

func (dj *delayedJobs) Len() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(*dj)
}

func (dj delayedJobs) Less(i, j int) bool {
	mu.RLock()
	defer mu.RUnlock()
	// We want Pop to give us the lowest, not higuest, priority so we use lower than here.
	return dj[i].priority < dj[j].priority
}

func (dj delayedJobs) Swap(i, j int) {
	mu.Lock()
	defer mu.Unlock()
	dj[i], dj[j] = dj[j], dj[i]
	dj[i].index = i
	dj[j].index = j
}

func (dj *delayedJobs) Push(x interface{}) {
	mu.Lock()
	defer mu.Unlock()
	n := len(*dj)
	item := x.(*job)
	item.index = n
	*dj = append(*dj, item)
}

func (dj *delayedJobs) Pop() interface{} {
	mu.Lock()
	defer mu.Unlock()
	old := *dj
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*dj = old[0 : n-1]
	return item
}

func (dj *delayedJobs) NextDeadLine() time.Time {
	mu.RLock()
	defer mu.RUnlock()
	next := *dj
	n := len(next)
	return next[n-1].Data.Deadline()
}

// update modifies the priority and value of an Item in the queue.
func (dj *delayedJobs) Update(item *job, data IJob, priority int64) {
	mu.Lock()
	defer mu.Unlock()
	heap.Remove(dj, item.index)
	item.Data = data
	item.priority = priority
	heap.Push(dj, item)
}

func (dj *delayedJobs) GetJobDetails(id bson.ObjectId) IJob {
	mu.RLock()
	defer mu.RUnlock()
	for _, job := range *dj {
		if job.Data.GetId() == id {
			return job.Data
		}
	}
	return nil
}

func (dj *delayedJobs) RemoveJob(id bson.ObjectId) {
	mu.RLock()
	defer mu.RUnlock()
	for index, job := range *dj {
		if job.Data.GetId() == id {
			tmp := *dj
			*dj = append(tmp[:index], tmp[index+1:]...)
		}
	}
}

func HandledelayedJobs(s *Scheduler) {
	c := time.Tick(200 * time.Millisecond)
	for now := range c {
		if s.jobs.Len() > 0 {
			nowUnix := now.UTC().Truncate(time.Second).Unix()
			deadline := s.jobs.NextDeadLine().UTC().Truncate(time.Second).Unix()
			if nowUnix > deadline {
				job := heap.Pop(s.jobs).(*job)
				go job.Data.Run()
			}
		}
	}
}
