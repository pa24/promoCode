package storage

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestApplyPromoCode_UsageLimitReached(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storageDB := &DB{DB: db}

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "reward", "max_uses", "uses"}).
		AddRow(1, 100, 5, 5)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, reward, max_uses, uses FROM promocodes WHERE code = $1 FOR UPDATE")).
		WithArgs("PROMO123").
		WillReturnRows(rows)

	err = ApplyPromoCode(storageDB, 1, "PROMO123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "promocode usage limit reached")
	assert.NoError(t, mock.ExpectationsWereMet())
}
