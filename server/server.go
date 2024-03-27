package server

import (
	"context"
	"net/http"
	"sfjgoapic/websocket"
	"time"
)

type Server struct {
	Addr   string
	Router http.Handler
	WsAddr string // WebSocket server address
}

func NewServer(addr string, router http.Handler, wsAddr string) *Server {
	return &Server{
		Addr:   addr,
		Router: router,
		WsAddr: wsAddr,
	}
}

func (s *Server) ListenAndServe() error {
	// Start ws server
	go websocket.StartWebSocketServer(s.WsAddr)

	// Start HTTP server
	return http.ListenAndServe(s.Addr, s.Router)
}

func (s *Server) Shutdown(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	server := &http.Server{Addr: s.Addr}

	err := server.Shutdown(ctx)
	if err != nil {
		return err
	}

	// Wait for all connections to be closed
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
