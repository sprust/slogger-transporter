package objects

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Tries   int
}

type QueueSettings struct {
	QueueName       string
	QueueWorkersNum int
}

type QueueInterface interface {
	GetSettings() (*QueueSettings, error)
	Publish(payload []byte) error
	Handle(job *Job) error
}

type Job struct {
	WorkerId int
	Payload  []byte
}
