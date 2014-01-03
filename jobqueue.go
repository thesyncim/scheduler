package scheduler

import (
	"container/heap"
	"labix.org/v2/mgo/bson"
	//"log"
	"sync"
	"time"
)

type jobQueue struct {
	jobs []*job
	mu   sync.RWMutex
}

func (jq *jobQueue) Len() int {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	return len(jq.jobs)
}

func (jq jobQueue) Less(i, j int) bool {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	// We want Pop to give us the lowest, not higuest, priority so we use lower than here.
	return jq.jobs[i].priority < jq.jobs[j].priority
}

func (jq jobQueue) Swap(i, j int) {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	jq.jobs[i], jq.jobs[j] = jq.jobs[j], jq.jobs[i]
	jq.jobs[i].index = i
	jq.jobs[j].index = j
}

func (jq *jobQueue) Push(x interface{}) {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	n := len(jq.jobs)
	item := x.(*job)
	item.index = n
	jq.jobs = append(jq.jobs, item)
}

func (jq *jobQueue) Pop() interface{} {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	old := jq.jobs
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	jq.jobs = old[0 : n-1]
	return item
}

func (jq *jobQueue) NextDeadLine() time.Time {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	//n := len(jq.jobs)
	return jq.jobs[0].Data.Deadline()
}

// update modifies the priority and value of an Item in the queue.
func (jq *jobQueue) Update(item *job, data Job, priority int64) {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	heap.Remove(jq, item.index)
	item.Data = data
	item.priority = priority
	heap.Push(jq, item)
}

func (jq *jobQueue) GetJobDetails(id bson.ObjectId) Job {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	for _, job := range jq.jobs {
		if job.Data.GetId() == id {
			return job.Data
		}
	}
	return nil
}

func (jq *jobQueue) RemoveJob(id bson.ObjectId) {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	for index, job := range jq.jobs {
		if job.Data.GetId() == id {
			tmp := jq.jobs
			jq.jobs = append(tmp[:index], tmp[index+1:]...)
		}
	}
}
