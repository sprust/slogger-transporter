package queue_trace_transporter

type Payload struct {
	Token   string `json:"tok"`
	Actions []Data `json:"acs"`
}

type Data struct {
	Type string `json:"tp"`
	Data string `json:"dt"`
}

type CreatingTrace struct {
	TraceId       string   `json:"tid"`
	ParentTraceId *string  `json:"pid"`
	Type          string   `json:"tp"`
	Status        string   `json:"st"`
	Tags          []string `json:"tgs"`
	Data          string   `json:"dt"`
	Duration      *float64 `json:"dur"`
	Memory        *float64 `json:"mem"`
	Cpu           *float64 `json:"cpu"`
	LoggedAt      string   `json:"lat"`
}

type UpdatingTrace struct {
	TraceId   string    `json:"tid"`
	Status    string    `json:"st"`
	Profiling *string   `json:"pr"`
	Tags      *[]string `json:"tgs"`
	Data      *string   `json:"dt"`
	Duration  *float64  `json:"dur"`
	Memory    *float64  `json:"mem"`
	Cpu       *float64  `json:"cpu"`
}
