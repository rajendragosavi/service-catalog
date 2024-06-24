package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	"github.com/rajendragosavi/service-catalog/pkg/db"
)

type Repository struct {
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Db: db}
}

// func (r Repository) Find(ctx context.Context, id int) (*model.ServiceCatalog, error) {
// 	obj := &model.ServiceCatalog{}
// 	query := "SELECT * FROM service WHERE id = $1 AND deleted_time IS NULL"

// 	err := r.Db.GetContext(ctx, &obj, query)
// 	return obj, db.HandleError(err)
// }

func (r *Repository) Create(ctx context.Context, obj *model.ServiceCatalog) error {
	query := `INSERT INTO service (service_name, description, status, creation_time, last_updated_time,deletion_time, versions,is_deleted)
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
	rows, err := r.Db.NamedQueryContext(ctx, query, params)
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

func (r *Repository) Get(ctx context.Context, serviceName string) (*model.ServiceCatalog, error) {
	obj := &model.ServiceCatalog{}
	query := "SELECT * FROM service WHERE service_name = $1 AND is_deleted IS FALSE"
	err := r.Db.GetContext(ctx, obj, query, serviceName)
	return obj, db.HandleError(err)
}
