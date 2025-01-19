package jobq

type Job interface {
	ID() string
}

type JobQueue struct {
	wait  chan int
	jobs  chan Job
	queue chan Job
}

func (jq JobQueue) Enqueue(jobs []Job) {
	if len(jobs) == 0 {
		return
	}
	jq.wait <- len(jobs)
	go func() {
		for _, j := range jobs {
			jq.queue <- j
		}
	}()
}

func (jq JobQueue) Jobs() <-chan Job {
	return jq.jobs
}

func (jq JobQueue) MarkJobDone() {
	jq.wait <- -1
}

func New() *JobQueue {
	queueCount := 0
	wait := make(chan int)
	jobs := make(chan Job)
	queue := make(chan Job)
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

	return &JobQueue{wait: wait, jobs: jobs, queue: queue}
}

type (
	SetProcessedFunc func(j Job)
	IsProcessedFunc  func(j Job) bool
)

func NewWithFuncs(isProcessed IsProcessedFunc, setProcessed SetProcessedFunc) *JobQueue {
	queueCount := 0
	wait := make(chan int)
	jobs := make(chan Job)
	queue := make(chan Job)

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
			if !isProcessed(j) {
				setProcessed(j)
				jobs <- j
			} else {
				wait <- -1
			}
		}

		close(jobs)
		close(wait)
	}()

	return &JobQueue{wait: wait, jobs: jobs, queue: queue}
}
