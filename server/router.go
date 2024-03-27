package server

import (
	"fmt"
	"net/http"
	handler "sfjgoapic/handlers"
	"sfjgoapic/websocket"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type XYRouter struct {
	router *mux.Router
	port   int
}

func NewRouter(h handler.HandlerInterface, port int) *XYRouter {
	r := mux.NewRouter()

	// Logging middleware to log incoming requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Incoming request:", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// CORS middleware configuration
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // Adjust with your frontend URL
	r.Use(handlers.CORS(headers, methods, origins))

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
