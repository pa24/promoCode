package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"promoCode/internal/models"
)

type FakePromoServiceAdmin struct {
	CreatePromoCodeFunc func(request models.CreatePromoRequest) error
}

func (f *FakePromoServiceAdmin) CreatePromoCode(request models.CreatePromoRequest) error {
	if f.CreatePromoCodeFunc != nil {
		return f.CreatePromoCodeFunc(request)
	}
	return nil
}

func (f *FakePromoServiceAdmin) ApplyPromoCode(request models.ApplyRequest) error {
	return nil
}

func TestAdminHandler_CreatePromoCode_CallsService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var called bool
	fakeService := &FakePromoServiceAdmin{
		CreatePromoCodeFunc: func(request models.CreatePromoRequest) error {
			called = true
			assert.Equal(t, "PROMO123", request.Code)
			assert.Equal(t, 100, request.Reward)
			assert.Equal(t, 5, request.MaxUses)
			return nil
		},
	}
	handler := NewAdminHandler(fakeService)

	form := url.Values{}
	form.Add("code", "PROMO123")
	form.Add("reward", "100")
	form.Add("max_uses", "5")

	req, err := http.NewRequest(http.MethodPost, "/admin/create", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePromoCode(c)

	assert.True(t, called, "Service CreatePromoCodeFunc should be called")
}

func TestAdminHandler_CreatePromoCode_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeService := &FakePromoServiceAdmin{}
	handler := NewAdminHandler(fakeService)

	// Отправляем запрос без обязательных данных.
	req, err := http.NewRequest(http.MethodPost, "/admin/create", strings.NewReader(""))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePromoCode(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
