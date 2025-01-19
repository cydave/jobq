package main

import (
	"log"
	"math/rand"
	"sync"

	"github.com/cydave/jobq"
	"github.com/google/uuid"
)

type Job struct {
	id string
}

func (j Job) ID() string {
	return j.id
}

func NewJob() Job {
	return Job{id: uuid.NewString()[:8]}
}

func worker(wg *sync.WaitGroup, jq *jobq.JobQueue) {
	workerID := uuid.NewString()[:8]
	defer wg.Done()
	for job := range jq.Jobs() {
		log.Printf("[worker:%s] Handling job #%s\n", workerID, job.ID())

		// Produce new jobs at random
		if rand.Intn(100)%2 == 0 {
			newJobsCount := rand.Intn(3)
			newJobs := make([]jobq.Job, newJobsCount)
			for i := 0; i < newJobsCount; i++ {
				j := NewJob()
				log.Printf("[worker:%s] Producing new job #%s\n", workerID, j.ID())
				newJobs[i] = j
			}
			if newJobsCount > 0 {
				// Produce duplicates to showcase the isProcessed and
				// setProcessed funcs.
				for i := 0; i < 3; i++ {
					newJobs = append(newJobs, job)
				}

				jq.Enqueue(newJobs)
			}
		}

		jq.MarkJobDone()
	}
}

func main() {
	processedJobs := make(map[string]interface{})
	isProcessed := func(j jobq.Job) bool {
		_, ok := processedJobs[j.ID()]
		if ok {
			log.Printf("[skipping] job #%s", j.ID())
		}
		return ok
	}
	setProcessed := func(j jobq.Job) {
		processedJobs[j.ID()] = struct{}{}
	}

	jq := jobq.NewWithFuncs(isProcessed, setProcessed)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(wg, jq)
	}

	jobs := make([]jobq.Job, 10)
	for i := 0; i < 10; i++ {
		jobs[i] = NewJob()
	}
	jq.Enqueue(jobs)
	wg.Wait()

	log.Println("all done!")
}
