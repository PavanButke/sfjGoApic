package handlers

import "net/http"

type HandlerInterface interface {
	SubmitJobHandler(w http.ResponseWriter, r *http.Request)
	GetJobsHandler(w http.ResponseWriter, r *http.Request)
}
