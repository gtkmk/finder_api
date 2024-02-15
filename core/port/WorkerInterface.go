package port

type WorkerTask interface {
	ExecuteWorker() error
}

type WorkerInterface interface {
	AddTask(task WorkerTask)
	Wait()
	HandleErrorsAsync(chan error)
	StartWorkers()
}
