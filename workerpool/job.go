package workerpool


type HandlerFunc func(interface{})

type Job struct {
	Input interface{}
	Handler HandlerFunc
}

type JobQueue chan Job
