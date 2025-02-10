package models

import "time"

type CreatePromoRequest struct {
	Code    string `form:"code" binding:"required"`
	Reward  int    `form:"reward" binding:"required"`
	MaxUses int    `form:"max_uses" binding:"required"`
}

type ApplyRequest struct {
	PlayerID int    `json:"player_id" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type PromoCode struct {
	Id        int
	Reward    int
	MaxUses   int
	Uses      int
	CreatedAt time.Time
}
