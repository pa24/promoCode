package service

import (
	"promoCode/internal/models"
	"promoCode/internal/storage"
)

type PromoService interface {
	CreatePromoCode(request models.CreatePromoRequest) error
	ApplyPromoCode(request models.ApplyRequest) error
}

type PromoServiceImpl struct {
	db *storage.DB
}

func NewPromoService(db *storage.DB) PromoService {
	return &PromoServiceImpl{db: db}
}

func (s *PromoServiceImpl) CreatePromoCode(request models.CreatePromoRequest) error {
	return storage.CreatePromoCode(s.db, request.Code, request.Reward, request.MaxUses)
}

func (s *PromoServiceImpl) ApplyPromoCode(request models.ApplyRequest) error {
	return storage.ApplyPromoCode(s.db, request.PlayerID, request.Code)
}
