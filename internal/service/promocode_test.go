package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"promoCode/internal/models"
	"promoCode/internal/storage"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromoService_CreatePromoCode(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storageDB := &storage.DB{DB: db}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO promocodes (code, reward, max_uses, created_at) VALUES ($1, $2, $3, $4)")).
		WithArgs("PROMO123", 100, 5, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	promoService := NewPromoService(storageDB)

	req := models.CreatePromoRequest{
		Code:    "PROMO123",
		Reward:  100,
		MaxUses: 5,
	}

	err = promoService.CreatePromoCode(req)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
