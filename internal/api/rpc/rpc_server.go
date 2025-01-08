package rpc

import (
	goridgeRpc "github.com/roadrunner-server/goridge/v3/pkg/rpc"
	"log/slog"
	"net"
	"net/rpc"
	"slogger-transporter/internal/api/rpc/ping_pong"
	"slogger-transporter/pkg/foundation/errs"
)

var functions = []any{
	&ping_pong.PingPong{},
}

type Server struct {
	rpcPort  string
	listener net.Listener
	closing  bool
}

func NewServer(rpcPort string) *Server {
	server := &Server{
		rpcPort: rpcPort,
	}

	return server
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		return errs.Err(err)
	}

	s.listener = listener

	for _, function := range functions {
		err = rpc.Register(function)

		if err != nil {
			slog.Error(err.Error())

			return errs.Err(err)
		}
	}

	slog.Info("Listening on port " + s.rpcPort)

	for {
		conn, err := s.listener.Accept()

		if s.closing == true {
			break
		}

		if err != nil {
			slog.Error("Error listening:", err.Error())

			continue
		}

		_ = conn

		go rpc.ServeCodec(goridgeRpc.NewCodec(conn))
	}

	return nil
}

func (s *Server) Close() error {
	slog.Warn("Closing rpc server...")

	s.closing = true

	return errs.Err(s.listener.Close())
}
