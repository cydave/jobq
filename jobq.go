package jobq

type Job interface {
	ID() string
}

type IJobQueue[T Job] interface {
	EnqueueSingle(T)
	Enqueue([]T)
	MarkJobDone()
	Jobs() <-chan T
}

type JobQueue[T Job] struct {
	wait  chan int
	jobs  chan T
	queue chan T
}

// Enqueue enqueues multiple jobs to be prolcessed.
func (jq *JobQueue[T]) Enqueue(jobs ...T) {
	if len(jobs) <= 0 {
		return
	}
	jq.wait <- len(jobs)
	go func() {
		for _, j := range jobs {
			jq.queue <- j
		}
	}()
}

// Jobs returns the job queue which workers can range over to prcess jobs.
func (jq *JobQueue[T]) Jobs() <-chan T {
	return jq.jobs
}

// MarkJobDone needs to be called after a job has been processed by a worker.
func (jq *JobQueue[T]) MarkJobDone() {
	jq.wait <- -1
}

func New[T Job]() *JobQueue[T] {
	queueCount := 0
	wait := make(chan int)
	jobs := make(chan T)
	queue := make(chan T)
	processed := map[string]interface{}{}

	go func() {
		for delta := range wait {
			queueCount += delta
			if queueCount == 0 {
				close(queue)
			}
		}
	}()

	go func() {
		for j := range queue {
			if _, ok := processed[j.ID()]; !ok {
				processed[j.ID()] = struct{}{}
				jobs <- j
			} else {
				wait <- -1
			}
		}

		close(jobs)
		close(wait)
	}()

	return &JobQueue[T]{wait: wait, jobs: jobs, queue: queue}
}
