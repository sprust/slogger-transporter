package ping_pong

type PingPong struct {
}

type PingPongArgs struct {
	Message string
}

type PingPongResult struct {
	Message string
}

func (p *PingPong) Pong(args *PingPongArgs, reply *PingPongResult) error {
	reply.Message = args.Message

	return nil
}
