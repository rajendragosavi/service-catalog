package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "service")
	repo := NewRepository(sqlxDB)

	// Create a new ServiceCatalog instance with the repository
	service := model.ServiceCatalog{
		Name:        "test_service",
		Description: "test_description",
		Status:      1,
		Versions:    pq.StringArray{"v1.0"},
		CreatedOn:   time.Now(),
	}

	// Set up mock expectations for the Create query
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO service \(service_name, description, status, creation_time, last_updated_time, deletion_time, versions, is_deleted\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING service_id;`).
		WithArgs(
			service.Name,
			service.Description,
			service.Status,
			service.CreatedOn,
			sql.NullTime{Valid: false}, // last_updated_time
			sql.NullTime{Valid: false}, // deletion_time
			pq.Array(service.Versions),
			service.IsDeleted,
		).
		WillReturnRows(sqlmock.NewRows([]string{"service_id"}).AddRow("513fa4a4-8bc2-4145-91f3-aca33c4eea8e"))
	mock.ExpectCommit()

	ctx := context.Background()
	err = repo.Create(ctx, &service)
	fmt.Printf("ERRROR - %+v \n", err)
	assert.NoError(t, err)
	assert.Equal(t, "513fa4a4-8bc2-4145-91f3-aca33c4eea8e", service.ID)

	//Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_Get(t *testing.T) {
	// Initialize a new mock DB
	// mockDB, mock, err := sqlmock.New()
	mockDB, mock, err := sqlxmock.Newx()
	require.NoError(t, err, "error opening mock database connection")
	defer mockDB.Close()

	// Create a sqlx.DB instance from the mock DB
	// sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Create a new repository with the mock DB
	repo := NewRepository(mockDB)

	// Expected query and arguments
	serviceName := "test_service"
	expectedQuery := "SELECT * FROM service WHERE service_name = $1 AND is_deleted IS FALSE"
	rows := sqlmock.NewRows([]string{"service_id", "service_name", "description", "status", "creation_time", "last_updated_time", "deletion_time", "versions", "is_deleted"}).
		AddRow("513fa4a4-8bc2-4145-91f3-aca33c4eea8e", "test_service", "test_description", 1, time.Now(), nil, nil, pq.StringArray{"v1.0"}, false)

	// Expectation for the query
	mock.ExpectQuery(expectedQuery).
		WithArgs(serviceName).
		WillReturnRows(rows)

	// Context for the test
	ctx := context.Background()

	// Call the Get method with the mocked data
	result, err := repo.Get(ctx, serviceName)

	// Verify expectations
	require.NoError(t, err, "expected no error from repository Get")
	assert.NotNil(t, result, "expected non-nil result")

	// Verify the fetched data
	expected := &model.ServiceCatalog{
		ID:          "513fa4a4-8bc2-4145-91f3-aca33c4eea8e",
		Name:        "test_service",
		Description: "test_description",
		Status:      1,
		Versions:    []string{"v1.0"},
		IsDeleted:   false,
		// Populate other fields as per your model definition
	}

	assert.Equal(t, expected, result, "expected fetched data to match")
}
