package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"promoCode/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type FakePromoService struct {
	ApplyPromoCodeFunc func(request models.ApplyRequest) error
}

func (f *FakePromoService) CreatePromoCode(request models.CreatePromoRequest) error {
	return nil
}

func (f *FakePromoService) ApplyPromoCode(request models.ApplyRequest) error {
	if f.ApplyPromoCodeFunc != nil {
		return f.ApplyPromoCodeFunc(request)
	}
	return nil
}

func TestAPIHandler_ApplyPromoCode_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подменяем метод ApplyPromoCode, чтобы он возвращал nil (успешное применение).
	fakeService := &FakePromoService{
		ApplyPromoCodeFunc: func(request models.ApplyRequest) error {
			return nil
		},
	}
	handler := NewAPIHandler(fakeService)

	reqBody, _ := json.Marshal(models.ApplyRequest{
		PlayerID: 1,
		Code:     "PROMO123",
	})
	req, err := http.NewRequest(http.MethodPost, "/api/apply", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ApplyPromoCode(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Promocode applied successfully", resp["message"])
}

func TestAPIHandler_ApplyPromoCode_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeService := &FakePromoService{}
	handler := NewAPIHandler(fakeService)

	req, err := http.NewRequest(http.MethodPost, "/api/apply", bytes.NewBufferString("{invalid"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ApplyPromoCode(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid JSON", resp["error"])
}

func TestAPIHandler_ApplyPromoCode_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeService := &FakePromoService{
		ApplyPromoCodeFunc: func(request models.ApplyRequest) error {
			return errors.New("promocode usage limit reached")
		},
	}
	handler := NewAPIHandler(fakeService)

	reqBody, _ := json.Marshal(models.ApplyRequest{
		PlayerID: 1,
		Code:     "PROMO123",
	})
	req, err := http.NewRequest(http.MethodPost, "/api/apply", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ApplyPromoCode(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "promocode usage limit reached", resp["error"])
}
