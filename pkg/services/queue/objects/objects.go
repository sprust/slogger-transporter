package objects

type RmqConfig struct {
	User string
	Pass string
	Host string
	Port string
}

type Message struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Tries   int    `json:"tries"`
}

type QueueSettings struct {
	QueueName       string
	QueueWorkersNum int
}

type QueueInterface interface {
	GetSettings() (*QueueSettings, error)
	Handle(job *Job) error
}

type Job struct {
	WorkerId int
	Payload  []byte
}
