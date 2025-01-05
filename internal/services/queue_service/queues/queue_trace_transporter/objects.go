package queue_trace_transporter

type Payload struct {
	Token string `json:"t"`
	Data  []Data `json:"d"`
}

type Data struct {
	Action string `json:"a"`
	Trace  string `json:"t"`
}

type CreatingTrace struct {
	TraceId       string   `json:"traceId"`
	ParentTraceId *string  `json:"parentTraceId"`
	Type          string   `json:"type"`
	Status        string   `json:"status"`
	Tags          []string `json:"tags"`
	Data          string   `json:"data"`
	Duration      *float64 `json:"duration"`
	Memory        *float64 `json:"memory"`
	Cpu           *float64 `json:"cpu"`
	LoggedAt      string   `json:"loggedAt"`
}

type UpdatingTrace struct {
	TraceId   string    `json:"traceId"`
	Status    string    `json:"status"`
	Profiling *string   `json:"profiling"`
	Tags      *[]string `json:"tags"`
	Data      *string   `json:"data"`
	Duration  *float64  `json:"duration"`
	Memory    *float64  `json:"memory"`
	Cpu       *float64  `json:"cpu"`
}
