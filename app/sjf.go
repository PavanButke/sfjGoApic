package app

import (
	"database/sql"
	"sfjgoapic/models"
	"sfjgoapic/websocket"
	"strconv"
)

type SJFScheduler struct {
	db *sql.DB
}

func NewSJFScheduler(db *sql.DB) *SJFScheduler {
	return &SJFScheduler{
		db: db,
	}
}

// AddJob adds a new job
func (s *SJFScheduler) AddJob(job models.Job) error {

	durationStr := job.Duration

	_, err := s.db.Exec("INSERT INTO job (name, duration, status) VALUES ($1, $2, $3)", job.Name, durationStr, job.Status)
	if err != nil {
		return err
	}

	// check status
	prevStatus, err := s.getPreviousStatus(job.ID)
	if err != nil {
		return err
	}
	prevJobStatus := models.JobStatus(prevStatus)

	if job.Status != prevJobStatus {
		// Broadcast job status update to WebSocket clients
		// Convert job.Status to string before concatenating
		msg := "Job status updated: " + job.Name + " - Status: " + string(job.Status)
		msgBytes := []byte(msg)

		websocket.BroadcastMessage(msgBytes)
	}

	return nil
}

func (s *SJFScheduler) UpdateJobStatus(jobID int, newStatus string) error {
	_, err := s.db.Exec("UPDATE job SET status = $1 WHERE id = $2", newStatus, jobID)
	if err != nil {
		return err
	}

	msg := []byte("Job status updated: ID - " + strconv.Itoa(jobID) + ", Status: " + newStatus)
	websocket.BroadcastMessage(msg)

	return nil
}

func (s *SJFScheduler) getPreviousStatus(jobID int) (string, error) {
	var prevStatus string
	err := s.db.QueryRow("SELECT status FROM job WHERE id = $1", jobID).Scan(&prevStatus)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return prevStatus, nil
}

func (s *SJFScheduler) Schedule() ([]models.Job, error) {
	rows, err := s.db.Query("SELECT id, name, duration, status FROM job  ORDER BY duration")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduledJobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Name, &job.Duration, &job.Status)
		if err != nil {
			return nil, err
		}

		scheduledJobs = append(scheduledJobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scheduledJobs, nil
}
