package repository

import (
	"context"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewRepository(db)

	tests := []struct {
		name    string
		s       Repository
		service model.ServiceCatalog
		mock    func()
		want    string
		wantErr bool
	}{
		{
			//When everything works as expected
			name:    "successfully inserted the service entry in the table",
			s:       s,
			service: model.ServiceCatalog{Name: "test-1", Description: "test description", Status: 1, CreatedOn: time.Now().UTC(), UpdatedOn: nil, DeletedOn: nil, Versions: pq.StringArray{"1.0"}, IsDeleted: false},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"service_id"}).AddRow("12345")
				mock.ExpectQuery("INSERT INTO service").WithArgs("test-1", "test description", 1, sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg(), pq.StringArray{"1.0"}, false).WillReturnRows(rows)
			},
			want: "12345",
		},
		// {
		// 	name:  "Empty Fields",
		// 	s:     s,
		// 	user: domain.User{
		// 		FirstName: "",
		// 		LastName:  "",
		// 		Username:  "username",
		// 		Password:  "password",
		// 	},
		// 	mock: func() {
		// 		rows := sqlxmock.NewRows([]string{"id"})
		// 		mock.ExpectQuery("INSERT INTO users").WithArgs("first_name", "last_name", "username", "password").WillReturnRows(rows)
		// 	},
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(context.Background(), &tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewRepository(db)
	creationTime := time.Now().UTC()
	// mock.ExpectBegin()
	// before we actually execute our api function, we need to expect required DB actions
	tests := []struct {
		name        string
		s           Repository
		serviceName string
		mock        func()
		want        *model.ServiceCatalog
		wantErr     bool
	}{
		{
			//When everything works as expected
			name:        "successfully selected service from service table",
			s:           s,
			serviceName: "test-service",
			mock: func() {
				rows := sqlxmock.NewRows([]string{"service_id", "service_name", "description", "status", "creation_time", "last_updated_time", "deletion_time", "versions", "is_deleted"}).
					AddRow("123456789", "test-1", "test-1-description", 1, creationTime, nil, nil, pq.StringArray{"1.0"}, false)
				mock.ExpectQuery("SELECT \\* FROM service WHERE service_name = \\$1 AND is_deleted IS FALSE").WithArgs("test-service").WillReturnRows(rows)
			},
			want: &model.ServiceCatalog{
				ID:          "123456789",
				Name:        "test-1",
				Description: "test-1-description",
				Status:      1,
				CreatedOn:   creationTime,
				UpdatedOn:   nil,
				DeletedOn:   nil,
				Versions:    pq.StringArray{"1.0"},
				IsDeleted:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Get(context.Background(), tt.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_List(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewRepository(db)
	creationTime := time.Now().UTC()
	serviceList := make([]*model.ServiceCatalog, 0)
	serviceList = append(serviceList, &model.ServiceCatalog{
		ID:          "123456789",
		Name:        "test-1",
		Description: "test-1-description",
		Status:      1,
		CreatedOn:   creationTime,
		UpdatedOn:   nil,
		DeletedOn:   nil,
		Versions:    pq.StringArray{"1.0"},
		IsDeleted:   false,
	}, &model.ServiceCatalog{
		ID:          "123456789",
		Name:        "test-2",
		Description: "test-3-description",
		Status:      1,
		CreatedOn:   creationTime,
		UpdatedOn:   nil,
		DeletedOn:   nil,
		Versions:    pq.StringArray{"1.0"},
		IsDeleted:   false,
	})
	// before we actually execute our api function, we need to expect required DB actions
	tests := []struct {
		name            string
		s               Repository
		totalNoServices int
		mock            func()
		want            []*model.ServiceCatalog
		wantErr         bool
	}{
		{
			//When everything works as expected
			name:            "successfully selected service from service table",
			s:               s,
			totalNoServices: 2,
			mock: func() {
				rows := sqlxmock.NewRows([]string{"service_id", "service_name", "description", "status", "creation_time", "last_updated_time", "deletion_time", "versions", "is_deleted"}).
					AddRow("123456789", "test-1", "test-1-description", 1, creationTime, nil, nil, pq.StringArray{"1.0"}, false).AddRow("123456789", "test-2", "test-2-description", 1, creationTime, nil, nil, pq.StringArray{"1.0"}, false)
				mock.ExpectQuery("SELECT \\* FROM service WHERE is_deleted IS FALSE").WillReturnRows(rows)
			},
			want: serviceList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.List(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && len((tt.want)) != 2 {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
