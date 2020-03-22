package asyncq

import "log"

//Dispatcher is a job dispatcher that dispatches job to
//worker job channel if it is available
type Dispatcher struct {
	Workerpool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
	End        chan bool
	Workers    []*Worker
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		Workerpool: pool,
		MaxWorkers: maxWorkers,
		JobQueue:   make(chan Job, 100),
	}
}

//Run fires up workers
func (d *Dispatcher) Run() {
	log.Println("[[asyncq/dispatcher]] Creating worker pool")
	workers := make([]*Worker, 0)
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.Workerpool)
		worker.Start()
		workers = append(workers, worker)
	}
	d.Workers = workers
	go d.dispatch()
}

//dispatch dispatches the job to worker job channel which is available
func (d *Dispatcher) dispatch() {
	for {
		select {
		//when job request is received
		case job := <-d.JobQueue:
			go func(job Job) {
				//obtain worker job channel if available
				//otherwise block
				jobChannel := <-d.Workerpool

				//dispatch the job to worker job channel
				jobChannel <- job
			}(job)
		case <-d.End:
			for _, w := range d.Workers {
				w.Stop() //stop workers one by one
			}
		}
	}
}
