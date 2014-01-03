package main

import (
	"fmt"
	"github.com/thesyncim/scheduler"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

func main() {
	s := scheduler.New()

	for i := 1; i <= 10; i++ {
		d := time.Now().Add(time.Duration(i) * time.Second)
		s.AddJob(newTask(d))

	}

	fmt.Println(s.GetJobs())

	select {}

}

func newTask(deadline time.Time) (t *task) {
	t = new(task)
	t.runAt = deadline
	t.id = bson.NewObjectId()
	return
}

type task struct {
	id    bson.ObjectId
	runAt time.Time
}

func (t *task) Deadline() time.Time {
	return t.runAt
}
func (t *task) Run() {
	log.Println("==>is running : ", t.GetId(), t.runAt)

}

func (t *task) GetId() bson.ObjectId {
	return t.id

}
