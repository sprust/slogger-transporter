package trace_transporter_service

type CreatingTrace struct {
	TraceId       string
	ParentTraceId *string
	Type          string
	Status        string
	Tags          []string
	Data          string
	Duration      *float64
	Memory        *float64
	Cpu           *float64
	LoggedAt      string
}

type UpdatingTrace struct {
	TraceId   string
	Status    string
	Profiling *string
	Tags      *[]string
	Data      *string
	Duration  *float64
	Memory    *float64
	Cpu       *float64
}
