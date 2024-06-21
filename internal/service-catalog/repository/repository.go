package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	"github.com/rajendragosavi/service-catalog/pkg/db"
)

type Repository struct {
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Db: db}
}

func (r Repository) Find(ctx context.Context, id int) (*model.ServiceCatalog, error) {
	obj := &model.ServiceCatalog{}
	query := "SELECT * FROM service WHERE id = $1 AND deleted_time IS NULL"

	err := r.Db.GetContext(ctx, &obj, query)
	return obj, db.HandleError(err)
}

func (r *Repository) Create(ctx context.Context, obj *model.ServiceCatalog) error {
	query := `INSERT INTO service (name, description, status, created_on, updated_on)
				VALUES (:name, :description, :status, :created_on, :updated_on) RETURNING id;`
	rows, err := r.Db.NamedQueryContext(ctx, query, obj)
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
