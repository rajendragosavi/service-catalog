package model

import (
	"time"

	"github.com/lib/pq"
)

type Status int

const (
	StatusPending Status = iota + 1
	StatusInProgress
	StatusDone
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending:
		return true
	case StatusInProgress:
		return true
	case StatusDone:
		return true
	}
	return false
}

type ServiceCatalog struct {
	ID          string `db:"service_id"`
	Name        string `db:"service_name"`
	Description string `db:"description"`
	Status      Status `db:"status"`
	//Versions    []string   `db:"versions"`
	Versions  pq.StringArray `db:"versions"`
	CreatedOn time.Time      `db:"creation_time"`
	UpdatedOn *time.Time     `db:"last_updated_time"`
	DeletedOn *time.Time     `db:"deletion_time"`
	IsDeleted bool           `db:"is_deleted"`
}
