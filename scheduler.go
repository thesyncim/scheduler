package scheduler

import (
	"container/heap"
	"labix.org/v2/mgo/bson"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	jobs *delayedJobs
	sync.RWMutex
}

func New() (s *Scheduler) {
	s = &Scheduler{}
	s.jobs = &delayedJobs{}
	heap.Init(s.jobs)
	log.Println("starting scheduler")
	go HandledelayedJobs(s)
	return
}

func (s *Scheduler) AddJob(ijob IJob) {
	j := &job{}
	j.Data = ijob
	j.priority = ijob.Deadline().UTC().Truncate(time.Second).Unix()
	heap.Push(s.jobs, j)
}

func (s *Scheduler) GetJobs() delayedJobs {
	return *(s.jobs)
}

func (s *Scheduler) RemoveJob(id bson.ObjectId) {
	s.jobs.RemoveJob(id)
}

func (s *Scheduler) GetJob(id bson.ObjectId) IJob {
	return s.jobs.GetJobDetails(id)
}

//func (s *Scheduler) UpdateJob(job IJob) {
//	job.GetId()
//	return s.jobs.Update(item, data, priority)
//}
