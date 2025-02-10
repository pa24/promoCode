package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"promoCode/internal/models"
	"time"
)

type DB struct {
	*sql.DB
}

func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		return nil, err
	}

	return &DB{db}, nil
}

func ApplyPromoCode(db *DB, playerID int, code string) error {
	tx, err := db.Begin()
	if err != nil {
		slog.Error("Failed to start transaction", "error", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("Failed to rollback transaction", "error", err)
		}
	}()

	promo, err := getPromoCode(tx, code)
	if err != nil {
		slog.Error("Failed to get promo code", "code", code, "error", err)
		return fmt.Errorf("failed to get promocode: %w", err)
	}

	if promo.Uses >= promo.MaxUses {
		slog.Warn("Promo code usage limit reached", "code", code, "uses", promo.Uses, "maxUses", promo.MaxUses)
		return errors.New("promocode usage limit reached")
	}

	used, err := hasPlayerUsedPromo(tx, playerID, promo.Id)
	if err != nil {
		slog.Error("Failed to check if player used promo code", "playerID", playerID, "code", code, "error", err)
		return fmt.Errorf("failed to check promocode usage: %w", err)
	}
	if used {
		slog.Warn("Player has already used this promo code", "playerID", playerID, "code", code)
		return errors.New("player has already used this promocode")
	}

	if err := applyReward(tx, playerID, promo); err != nil {
		slog.Error("Failed to apply reward", "playerID", playerID, "code", code, "error", err)
		return fmt.Errorf("failed to apply reward: %w", err)
	}

	if err := incrementPromoUsage(tx, promo.Id); err != nil {
		slog.Error("Failed to increment promo code usage", "code", code, "error", err)
		return fmt.Errorf("failed to increment promocode usage: %w", err)
	}

	return tx.Commit()

}

func CreatePromoCode(db *DB, code string, reward, maxUses int) error {
	slog.Info("Creating promo code", "code", code, "reward", reward, "maxUses", maxUses)
	_, err := db.Exec("INSERT INTO promocodes (code, reward, max_uses, created_at) VALUES ($1, $2, $3, $4)",
		code, reward, maxUses, time.Now())
	if err != nil {
		slog.Error("Failed to create promo code", "code", code, "error", err)
		return err
	}
	return nil
}

// todo check
func getPromoCode(tx *sql.Tx, code string) (models.PromoCode, error) {
	var promo models.PromoCode
	err := tx.QueryRow("SELECT id, reward, max_uses, uses FROM promocodes WHERE code = $1 FOR UPDATE", code).Scan(&promo.Id, &promo.Reward, &promo.MaxUses, &promo.Uses)
	if errors.Is(sql.ErrNoRows, err) {
		slog.Warn("Promo code not found", "code", code)
		return models.PromoCode{}, errors.New("promocode not found")
	} else if err != nil {
		slog.Error("Failed to get promo code from DB", "code", code, "error", err)
		return models.PromoCode{}, fmt.Errorf("failed to fetch promocode from DB: %w", err)
	}
	return promo, nil
}

func hasPlayerUsedPromo(tx *sql.Tx, playerID, promoID int) (bool, error) {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS (SELECT 1 FROM rewards WHERE player_id = $1 AND promocode_id = $2)", playerID, promoID).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check player used promo code", "playerID", playerID, "promoID", promoID, "error", err)
		return false, fmt.Errorf("failed to check if player used promocode: %w", err)
	}
	return exists, nil
}

func applyReward(tx *sql.Tx, playerID int, promo models.PromoCode) error {
	_, err := tx.Exec("INSERT INTO rewards (player_id, promocode_id, reward, applied_at) VALUES ($1, $2, $3, $4)",
		playerID, promo.Id, promo.Reward, time.Now())
	if err != nil {
		slog.Error("Failed to apply reward", "playerID", playerID, "promoID", promo.Id, "error", err)
		return fmt.Errorf("failed to apply reward: %w", err)
	}
	return nil
}

func incrementPromoUsage(tx *sql.Tx, promoID int) error {
	_, err := tx.Exec("UPDATE promocodes SET uses = uses + 1 WHERE id = $1", promoID)
	if err != nil {
		slog.Error("Failed to increment promo code usage", "promoID", promoID, "error", err)
		return fmt.Errorf("failed to update promocode usage: %w", err)
	}
	return nil
}
