package repository

import (
	"context"
	"database/sql"

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
}
type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, obj *model.ServiceCatalog) (string, error) {
	var service_id string
	row := r.db.QueryRow("INSERT INTO service (service_name, description, status, creation_time, last_updated_time,deletion_time, versions, is_deleted) VALUES ($1, $2, $3, $4 ,$5, $6 , $7, $8) RETURNING service_id",
		obj.Name, obj.Description, obj.Status, obj.CreatedOn, obj.UpdatedOn, obj.DeletedOn, obj.Versions, obj.IsDeleted)
	if err := row.Scan(&service_id); err != nil {
		return "", db.HandleError(err)
	}
	obj.ID = service_id
	return service_id, nil
}

func (r *repository) CreatewithCommit(ctx context.Context, obj *model.ServiceCatalog) error {
	query := `INSERT INTO service (service_name, description, status, creation_time, last_updated_time,deletion_time, versions, is_deleted)
				VALUES (:service_name, :description, :status, :creation_time, :last_updated_time , :deletion_time , :versions, :is_deleted) RETURNING service_id;`

	// Use a map to bind parameters, including converting Versions to pq.Array
	params := map[string]interface{}{
		"service_name":      obj.Name,
		"description":       obj.Description,
		"status":            obj.Status,
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
