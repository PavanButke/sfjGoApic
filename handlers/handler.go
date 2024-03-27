package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sfjgoapic/app"
	"sfjgoapic/models"
	"sfjgoapic/websocket"
)

type Handler struct {
	Scheduler *app.SJFScheduler
}

func NewHandler(scheduler *app.SJFScheduler) *Handler {
	return &Handler{Scheduler: scheduler}
}

func (h *Handler) SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload into a Job struct
	var job models.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if job.Duration <= "0" {
		http.Error(w, "Duration must be a positive integer", http.StatusBadRequest)
		return
	}

	if job.Status == "" {
		job.Status = models.Pending
	}

	err = h.Scheduler.AddJob(job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  job status update to ws clients
	msg := []byte("New job submitted: " + job.Name)
	websocket.BroadcastMessage(msg)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.Scheduler.Schedule()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving jobs: %v", err), http.StatusInternalServerError)
		return
	}

	if len(jobs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Marshal the list of jobs into JSON format
	response, err := json.Marshal(jobs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling jobs to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}
