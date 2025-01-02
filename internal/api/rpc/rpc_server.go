package rpc

import (
	goridgeRpc "github.com/roadrunner-server/goridge/v3/pkg/rpc"
	"log/slog"
	"net"
	"net/rpc"
	"slogger-transporter/internal/api/rpc/ping_pong"
	"slogger-transporter/internal/app"
)

var functions = []any{
	&ping_pong.PingPong{},
}

type Server struct {
	app      *app.App
	rpcPort  string
	listener net.Listener
	closing  bool
}

func NewServer(app *app.App, rpcPort string) *Server {
	server := &Server{
		app:     app,
		rpcPort: rpcPort,
	}

	app.AddCloseListener(server)

	return server
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		return err
	}

	s.listener = listener

	for _, function := range functions {
		err := rpc.Register(function)

		if err != nil {
			slog.Error(err.Error())

			return err
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
	s.closing = true

	return s.listener.Close()
}