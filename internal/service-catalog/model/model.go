package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	UserID   uuid.UUID `db:"user_id" json:"user_id"`
	UserName string    `db:"user_name" json:"user_name"`
	BUName   string    `db:"bu_name" json:"bu_name"`
}

// UserServiceAccess represents the access control between users and services
type UserServiceAccess struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	ServiceID uuid.UUID `db:"service_id" json:"service_id"`
}

type ServiceCatalog struct {
	ID          string `db:"service_id"`
	Name        string `db:"service_name"`
	Description string `db:"description"`
	Status      int    `db:"status"`
	//Versions    []string   `db:"versions"`
	CreatedOn time.Time      `db:"creation_time"`
	UpdatedOn *time.Time     `db:"last_updated_time"`
	DeletedOn *time.Time     `db:"deletion_time"`
	Versions  pq.StringArray `db:"versions"`
	IsDeleted bool           `db:"is_deleted"`
}
