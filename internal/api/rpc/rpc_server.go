package rpc

import (
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
	for _, function := range functions {
		err := rpc.Register(function)

		if err != nil {
			slog.Error(err.Error())

			return err
		}
	}

	listener, err := net.Listen("tcp", ":"+s.rpcPort)

	if err != nil {
		slog.Error("Error listening:", err.Error())

		return err
	}

	s.listener = listener

	slog.Info("Listening on port " + s.rpcPort)

	rpc.Accept(s.listener)

	return nil
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}
