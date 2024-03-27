package models

type JobStatus string

const (
	Pending   JobStatus = "Pending"
	Running   JobStatus = "Running"
	Completed JobStatus = "Completed"
)

type Job struct {
	ID       int       `db:"id,omitempty" json:"id,omitempty"`
	Name     string    `db:"name,omitempty" json:"name,omitempty"`
	Duration string    `db:"duration,omitempty" json:"duration,omitempty"`
	Status   JobStatus `db:"status,omitempty" json:"status,omitempty"`
}
