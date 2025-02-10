package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"promoCode/internal/models"
	"promoCode/internal/service"
)

type APIHandler struct {
	PromoService service.PromoService
}

func NewAPIHandler(promoService service.PromoService) *APIHandler {
	return &APIHandler{PromoService: promoService}
}

func (a *APIHandler) ApplyPromoCode(c *gin.Context) {
	slog.Info("Received request to apply promo code")
	var req models.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	slog.Info("Request data", "playerID", req.PlayerID, "code", req.Code)

	if req.PlayerID <= 0 || req.Code == "" {
		slog.Warn("Invalid request data", "playerID", req.PlayerID, "code", req.Code)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player_id or code"})
		return
	}

	if err := a.PromoService.ApplyPromoCode(req); err != nil {
		slog.Error("Failed to apply promo code", "playerID", req.PlayerID, "code", req.Code, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Promo code applied successfully", "playerID", req.PlayerID, "code", req.Code)
	c.JSON(http.StatusOK, gin.H{"message": "Promocode applied successfully"})
}
