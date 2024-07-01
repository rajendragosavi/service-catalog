package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	"github.com/rajendragosavi/service-catalog/pkg/db"
)

type Repository interface {
	Create(ctx context.Context, obj *model.ServiceCatalog) (string, error)
	Get(ctx context.Context, serviceName string) (*model.ServiceCatalog, error)
	List(ctx context.Context) ([]model.ServiceCatalog, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	CreatewithCommit(ctx context.Context, obj *model.ServiceCatalog) error
	// methods on user and access table
	CreateUser(ctx context.Context, user *model.User) error
	GrantAccess(ctx context.Context, userID, serviceID uuid.UUID) error
	CheckUserExists(ctx context.Context, userID string) (bool, error)
	CheckAccess(ctx context.Context, userID, serviceID uuid.UUID) (bool, error)
	ListServicesForUser(ctx context.Context, userID uuid.UUID) ([]model.ServiceCatalog, error)
	ListUsersWithAccess(ctx context.Context, serviceID uuid.UUID) ([]model.User, error)
}
type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, obj *model.ServiceCatalog) (string, error) {
	var service_id string
	row := r.db.QueryRow("INSERT INTO service (service_name, description, creation_time, last_updated_time,deletion_time, versions, is_deleted) VALUES ($1, $2, $3, $4 ,$5, $6 , $7) RETURNING service_id",
		obj.Name, obj.Description, obj.CreatedOn, obj.UpdatedOn, obj.DeletedOn, obj.Versions, obj.IsDeleted)
	if err := row.Scan(&service_id); err != nil {
		return "", db.HandleError(err)
	}
	obj.ID = service_id
	return service_id, nil
}

func (r *repository) CreatewithCommit(ctx context.Context, obj *model.ServiceCatalog) error {
	query := `INSERT INTO service (service_name, description, creation_time, last_updated_time,deletion_time, versions, is_deleted)
				VALUES (:service_name, :description, :creation_time, :last_updated_time , :deletion_time , :versions, :is_deleted) RETURNING service_id;`

	// Use a map to bind parameters, including converting Versions to pq.Array
	params := map[string]interface{}{
		"service_name":      obj.Name,
		"description":       obj.Description,
		"creation_time":     obj.CreatedOn,
		"last_updated_time": obj.UpdatedOn,
		"deletion_time":     obj.DeletedOn,
		"versions":          pq.Array(obj.Versions),
		"is_deleted":        obj.IsDeleted,
	}

	// obj.Versions = pq.Array(obj.Versions).([]string)
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return db.HandleError(err)
	}

	for rows.Next() {
		err = rows.StructScan(obj)
		if err != nil {
			return db.HandleError(err)
		}
	}
	return db.HandleError(err)
}

func (r *repository) Get(ctx context.Context, serviceName string) (*model.ServiceCatalog, error) {
	obj := &model.ServiceCatalog{}
	query := "SELECT * FROM service WHERE service_name = $1 AND is_deleted IS FALSE"
	err := r.db.GetContext(ctx, obj, query, serviceName)
	return obj, db.HandleError(err)
}

func (r *repository) List(ctx context.Context) ([]model.ServiceCatalog, error) {
	obj := make([]model.ServiceCatalog, 0)
	query := "SELECT * FROM service WHERE is_deleted IS FALSE"
	err := r.db.SelectContext(ctx, &obj, query)
	return obj, db.HandleError(err)
}

func (r *repository) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, opts)
}

/// *******************

func (r *repository) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (user_id, user_name, bu_name) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, user.UserID, user.UserName, user.BUName)
	return err
}

func (r *repository) CheckUserExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE user_id = $1)`
	err := r.db.GetContext(ctx, &exists, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}
func (r *repository) GrantAccess(ctx context.Context, userID, serviceID uuid.UUID) error {
	query := `INSERT INTO user_service_access (user_id, service_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, userID, serviceID)
	return err
}

func (r *repository) CheckAccess(ctx context.Context, userID, serviceID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM user_service_access WHERE user_id = $1 AND service_id = $2`
	var count int
	err := r.db.GetContext(ctx, &count, query, userID, serviceID)
	return count > 0, err
}

func (r *repository) ListServicesForUser(ctx context.Context, userID uuid.UUID) ([]model.ServiceCatalog, error) {
	var services []model.ServiceCatalog
	query := `SELECT s.* 
              FROM user_service_access usa
              JOIN service s ON usa.service_id = s.service_id
              WHERE usa.user_id = $1`
	err := r.db.SelectContext(ctx, &services, query, userID)
	return services, err
}

func (r *repository) ListUsersWithAccess(ctx context.Context, serviceID uuid.UUID) ([]model.User, error) {
	var users []model.User
	query := `SELECT u.* 
              FROM user_service_access usa
              JOIN users u ON usa.user_id = u.user_id
              WHERE usa.service_id = $1`
	err := r.db.SelectContext(ctx, &users, query, serviceID)
	return users, err
}
