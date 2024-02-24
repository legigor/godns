package dns

import (
	"context"
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	logger *slog.Logger
	ctx    context.Context
}

func NewServer(logger *slog.Logger, ctx context.Context) *Server {
	return &Server{
		logger: logger.With("context", "server"),
		ctx:    ctx,
	}
}

func (srv *Server) Start() (*net.UDPAddr, error) {

	addr := net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, fmt.Errorf("error opening the connection, %w", err)
	}

	go func(conn *net.UDPConn) {
		defer func(conn *net.UDPConn) {
			_ = conn.Close()
		}(conn)

		buf := make([]byte, 1024)
		for {
			select {
			case <-srv.ctx.Done():
				srv.logger.Info("stop listening")
				return
			default:
				n, addr, err := conn.ReadFromUDP(buf)
				if err != nil {
					srv.logger.Error("error reading data from UDP", "err", err)
					continue
				}

				request := string(buf[:n])
				srv.logger.Info("got a request", "message", request, "client", addr)

				_, err = conn.WriteToUDP([]byte("Hello, Client!"), addr)
				if err != nil {
					srv.logger.Error("error writing to UDP", "error", err)
				}
			}
		}

	}(conn)

	return conn.LocalAddr().(*net.UDPAddr), nil
}
