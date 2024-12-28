package ping_pong

import "log/slog"

type PingPong struct {
}

type PingPongArgs struct {
	Message string
}

type PingPongResult struct {
	Message string
}

func (p *PingPong) Ping(args *PingPongArgs, reply *PingPongResult) error {
	reply.Message = args.Message

	go slog.Info("Ping: " + args.Message)

	return nil
}
