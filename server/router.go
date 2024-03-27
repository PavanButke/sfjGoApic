// server/router.go

package server

import (
	handler "sfjgoapic/handlers"
	"sfjgoapic/websocket"

	"github.com/gorilla/mux"
)

// YourRouter represents the server router
type XYRouter struct {
	router *mux.Router
	port   int
}

func NewRouter(h handler.HandlerInterface, port int) *XYRouter {
	r := mux.NewRouter()

	r.HandleFunc("/jobs", h.SubmitJobHandler).Methods("POST")
	r.HandleFunc("/jobs", h.GetJobsHandler).Methods("GET")
		// Register WebSocket endpoint
	r.HandleFunc("/ws", websocket.HandleWebSocket)

	return &XYRouter{router: r, port: port}
}

func (r *XYRouter) Router() *mux.Router {
	return r.router
}

func (r *XYRouter) Port() int {
	return r.port
}
