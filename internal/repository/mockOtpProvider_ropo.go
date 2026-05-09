package repository

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"ROP_Backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	OTPCooldownMinutes = 2
	TokenLen           = 2
	PinLen             = 6
	refNoLen           = 6
	OTPExpiredTime     = 10 * time.Minute
)

type MockOTPRepository struct {
	db *gorm.DB
}

func generatePIN(length int) string {
	digits := "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

func generateRefNo(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func generateToken(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func NewMockOTPRepository(db *gorm.DB) *MockOTPRepository {
	return &MockOTPRepository{db: db}
}

func (r *MockOTPRepository) CreateOTP(ctx context.Context, tel string) (*models.MockOTP, error) {
	// ไม่ต้องมีเพราะprovider ให้ยิงได้ตลอด ไม่ว่าเบอร์ซ้ำor not
	now := time.Now()
	// cooldownTime := now.Add(-OTPCooldownMinutes * time.Minute)

	// var existing models.MockOTP

	// err := r.db.Where("tel = ? AND created_at >= ? AND status = ?",
	// 	tel, cooldownTime, models.StatusPending).
	// 	Order("created_at DESC").
	// 	First(&existing).Error

	// if err == nil {
	// 	return nil, errors.New("OTP already requested recently, please wait")
	// }

	// if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, err
	// }

	token := generateToken(TokenLen)
	refNo := generateRefNo(refNoLen)
	pin := generatePIN(PinLen)

	otp := &models.MockOTP{
		Tel:       tel,
		Token:     token,
		RefNo:     refNo,
		Pin:       pin,
		Status:    models.StatusPending,
		ExpiresAt: now.Add(OTPExpiredTime),
	}

	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

type VerifyReason string

const (
	ReasonNone         VerifyReason = ""
	ReasonInvalidToken VerifyReason = "invalid_token"
	ReasonInvalidPIN   VerifyReason = "invalid_pin"
	ReasonExpired      VerifyReason = "expired"
	ReasonAlreadyUsed  VerifyReason = "already_used"
)

func (r *MockOTPRepository) VerifyOTP(ctx context.Context, token string, pin string) (bool, VerifyReason, error) {

	var otp models.MockOTP

	tx := r.db.WithContext(ctx).Begin()

	// mimic cronjob
	if err := tx.
		Model(&models.MockOTP{}).
		Where("status = ? AND expires_at <= ?", models.StatusPending, time.Now()).
		Update("status", models.StatusExpired).Error; err != nil {

		tx.Rollback()
		return false, ReasonNone, err
	}

	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("token = ?", token).
		First(&otp).Error; err != nil {

		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ReasonInvalidToken, nil
		}
		return false, ReasonNone, err
	}

	if otp.Status == models.StatusVerified {
		tx.Rollback()
		return false, ReasonAlreadyUsed, nil
	}

	if otp.Status == models.StatusExpired || time.Now().After(otp.ExpiresAt) {
		tx.Model(&otp).Update("status", models.StatusExpired)
		tx.Commit()
		return false, ReasonExpired, nil
	}

	if otp.Pin != pin {
		tx.Rollback()
		return false, ReasonInvalidPIN, nil
	}

	if err := tx.Model(&otp).
		Where("status = ?", models.StatusPending).
		Update("status", models.StatusVerified).Error; err != nil {

		tx.Rollback()
		return false, ReasonNone, err
	}

	if err := tx.Commit().Error; err != nil {
		return false, ReasonNone, err
	}

	return true, ReasonNone, nil
}
