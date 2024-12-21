package repository_test

import (
	"errors"
	"project/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func TestGetTotalRevenueByStatus(t *testing.T) {

	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	log := zap.NewNop()

	repo := repository.NewRevenueRepository(db, log)

	mockResults := sqlmock.NewRows([]string{"status_payment", "revenue"}).
		AddRow("Completed", 100.50).
		AddRow("Cancelled", 50.25).
		AddRow("In Process", 75.10)

	mock.ExpectQuery(`SELECT status_payment, SUM\(total\) as revenue`).
		WillReturnRows(mockResults)

	response, err := repo.GetTotalRevenueByStatus()

	assert.NoError(t, err)
	assert.NotNil(t, response)

	totalRevenue := response["total_revenue"].(float64)
	revenueMap := response["by_status"].(map[string]float64)

	assert.Equal(t, 225.85, totalRevenue)
	assert.Equal(t, map[string]float64{
		"Completed":  100.50,
		"Cancelled":  50.25,
		"In Process": 75.10,
	}, revenueMap)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTotalRevenueByStatus_DBError(t *testing.T) {

	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	log := zap.NewNop()

	repo := repository.NewRevenueRepository(db, log)

	mock.ExpectQuery(`SELECT status_payment, SUM\(total\) as revenue`).
		WillReturnError(errors.New("database error"))

	response, err := repo.GetTotalRevenueByStatus()

	assert.Error(t, err)
	assert.Nil(t, response)

	assert.NoError(t, mock.ExpectationsWereMet())
}
