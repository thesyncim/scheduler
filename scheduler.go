package scheduler

import (
	"container/heap"
	"labix.org/v2/mgo/bson"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	jobs *jobQueue
	wait sync.WaitGroup
}

func handlejobQueue(s *Scheduler) {
	c := time.Tick(200 * time.Millisecond)
	for now := range c {

		if s.jobs.Len() > 0 {
			log.Println("tick", s.jobs.NextDeadLine())
			nowUnix := ToUnixTimestamp(now)
			deadline := ToUnixTimestamp(s.jobs.NextDeadLine())
			if nowUnix > deadline {
				job := heap.Pop(s.jobs).(*job)
				go job.Data.Run()
			}
		}
	}
}

func New() (s *Scheduler) {
	s = &Scheduler{}
	s.jobs = &jobQueue{}
	heap.Init(s.jobs)
	log.Println("starting scheduler")

	go handlejobQueue(s)

	return
}

func (s *Scheduler) AddJob(ijob Job) {
	j := &job{}
	j.Data = ijob
	j.priority = ToUnixTimestamp(ijob.Deadline())
	heap.Push(s.jobs, j)
}

func (s *Scheduler) GetJobs() jobQueue {
	return *(s.jobs)
}

func (s *Scheduler) RemoveJob(id bson.ObjectId) {
	s.jobs.RemoveJob(id)
}

func (s *Scheduler) GetJob(id bson.ObjectId) Job {
	return s.jobs.GetJobDetails(id)
}

//func (s *Scheduler) UpdateJob(job Job) {
//	job.GetId()
//	return s.jobs.Update(item, data, priority)
//}
