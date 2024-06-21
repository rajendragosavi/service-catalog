package model

import "time"

type Status int

const (
	StatusPending Status = iota + 1
	StatusInProgress
	StatusDone
)

type ServiceCatalog struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Status      Status     `db:"status"`
	Versions    []string   `db:"versions"`
	CreatedOn   time.Time  `db:"created_on"`
	UpdatedOn   *time.Time `db:"updated_on"`
	DeletedOn   *time.Time `db:"deleted_on"`
}
