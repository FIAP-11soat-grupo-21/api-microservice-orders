package data_source

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"microservice/internal/adapters/daos"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		sqlDB.Close()
	}

	return db, mock, cleanup
}

func TestGormOrderStatusDataSource_FindAll_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("status-1", "Pending").
		AddRow("status-2", "Paid")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "order_status"`)).
		WillReturnRows(rows)

	ds := &GormOrderStatusDataSource{db: db}

	result, err := ds.FindAll()
	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("FindAll() len = %d, want 2", len(result))
	}

	expected := daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}
	if result[0] != expected {
		t.Errorf("unexpected result[0]: %+v", result[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sql expectations: %v", err)
	}
}

func TestGormOrderStatusDataSource_FindAll_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "order_status"`)).
		WillReturnError(errors.New("db error"))

	ds := &GormOrderStatusDataSource{db: db}

	_, err := ds.FindAll()
	if err == nil {
		t.Error("FindAll() expected error, got nil")
	}
}

func TestGormOrderStatusDataSource_FindByID_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("status-1", "Pending")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "order_status" WHERE id = $1 ORDER BY "order_status"."id" LIMIT $2`)).
		WithArgs("status-1", sqlmock.AnyArg()).
		WillReturnRows(rows)

	ds := &GormOrderStatusDataSource{db: db}

	status, err := ds.FindByID("status-1")
	if err != nil {
		t.Fatalf("FindByID() unexpected error: %v", err)
	}

	expected := daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}
	if status != expected {
		t.Errorf("FindByID() = %+v, want %+v", status, expected)
	}
}

func TestGormOrderStatusDataSource_FindByID_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "order_status" WHERE id = $1 ORDER BY "order_status"."id" LIMIT $2`)).
		WithArgs("status-1", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	ds := &GormOrderStatusDataSource{db: db}

	_, err := ds.FindByID("status-1")
	if err == nil {
		t.Error("FindByID() expected error, got nil")
	}
}

func TestGormOrderStatusDataSource_FindByName_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("status-1", "Pending")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "order_status" WHERE name = $1 ORDER BY "order_status"."id" LIMIT $2`)).
		WithArgs("Pending", sqlmock.AnyArg()).
		WillReturnRows(rows)

	ds := &GormOrderStatusDataSource{db: db}

	status, err := ds.FindByName("Pending")
	if err != nil {
		t.Fatalf("FindByName() unexpected error: %v", err)
	}

	expected := daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}
	if status != expected {
		t.Errorf("FindByName() = %+v, want %+v", status, expected)
	}
}

func TestGormOrderStatusDataSource_FindByName_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "order_status" WHERE name = $1 ORDER BY "order_status"."id" LIMIT $2`)).
		WithArgs("Pending", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	ds := &GormOrderStatusDataSource{db: db}

	_, err := ds.FindByName("Pending")
	if err == nil {
		t.Error("FindByName() expected error, got nil")
	}
}
