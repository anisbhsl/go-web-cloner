package asyncq

//Job defines interface for any async JOB
type Job interface{
	Perform()
}

