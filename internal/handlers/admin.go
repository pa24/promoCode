package handlers

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log/slog"
	"net/http"
	"promoCode/internal/models"
	"promoCode/internal/service"
	"strings"
)

type AdminHandler struct {
	PromoService service.PromoService
}

func NewAdminHandler(promoService service.PromoService) *AdminHandler {
	return &AdminHandler{PromoService: promoService}
}

func (a *AdminHandler) AdminPage(c *gin.Context) {

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		slog.Error("Failed to load template", "error", err)
		c.String(http.StatusInternalServerError, "Error loading template")
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		slog.Error("Failed to render template", "error", err)
		c.String(http.StatusInternalServerError, "Error rendering template")
	}
	slog.Info("Admin page loaded successfully", "path", c.Request.URL.Path)
}

func (a *AdminHandler) CreatePromoCode(c *gin.Context) {
	var request models.CreatePromoRequest
	if err := c.ShouldBind(&request); err != nil {
		slog.Warn("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.PromoService.CreatePromoCode(request); err != nil {
		slog.Error("Failed to create promo code", "error", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Promo code already exist"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("Promo code created successfully", "code", request.Code)
		c.JSON(http.StatusOK, gin.H{"message": "Promo code created successfully"})

	}
}
