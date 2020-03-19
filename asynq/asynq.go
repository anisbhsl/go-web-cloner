package asyncq

import "log"

//JobQueue is a buffered Job channel where job requests will be sent
var JobQueue = make(chan Job)

//Worker executes the Job
type Worker struct{
	Workerpool chan chan Job
	JobChannel chan Job
	StopWorker chan bool
}

func NewWorker(workerPool chan chan Job) *Worker{
	log.Println("[[asyncq/worker]] Created a new worker")
	return &Worker{
		Workerpool: workerPool,
		JobChannel: make(chan Job),
		StopWorker: make(chan bool),
	}
}

//Start fires up the worker and keeps listening for any incoming job requests
func (w Worker) Start(){
	log.Println("[[asyncq/worker]] Started listening for incoming jobs")
	go func(){
		for{
			//Register a worker into job channel
			w.Workerpool<-w.JobChannel

			select{
			case job:= <-w.JobChannel:
				//whenever a new job arrives in job channel start execution
				log.Println("[[38:async]] new job received!")
				job.Perform() //perform a job
			case <-w.StopWorker:
				return //stop the worker when signal is received!
			}
		}
	}()
}


//Stop signals to stop the running worker
func (w Worker) Stop(){
	go func(){
		w.StopWorker<-true
	}()
}