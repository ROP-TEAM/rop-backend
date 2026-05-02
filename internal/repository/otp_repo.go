package repository

import (
	"ROP_Backend/internal/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	MaxRequestByUser        = 5
	UserRequestWindowPeriod = 1 * time.Hour
)

var (
	ErrReachMaxAttempt   = errors.New("verification reach the maximum attempt")
	ErrExpiredOTPRequest = errors.New("this otpRequest already expired")
	ErrUsedOTPRequest    = errors.New("tel already be verified by otp pin")
)

type OTPRepository struct {
	db *gorm.DB
}

type CreateOTPRequest struct {
	Tel    string
	Token  string
	UserId uint
	RefNo  string
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) WithTx(tx *gorm.DB) *OTPRepository {
	return &OTPRepository{db: tx}
}

func (r *OTPRepository) CreateOTPRequest(ctx context.Context, req *CreateOTPRequest) error {
	otp := models.OtpRequest{
		Tel:         req.Tel,
		Token:       req.Token,
		UserID:      req.UserId,
		RefNo:       req.RefNo,
		MaxAttempts: 5,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		IsUsed:      false,
		Attempts:    0,
	}

	if err := r.db.WithContext(ctx).Create(&otp).Error; err != nil {
		return err
	}

	return nil
}

func (r *OTPRepository) IsOTPRequestMutable(ctx context.Context, tel string, ref string) (string, error) {
	var otp models.OtpRequest

	err := r.db.WithContext(ctx).
		Where("tel = ? AND ref_no = ?", tel, ref).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return "", err
	}

	isExpired := time.Now().After(otp.ExpiresAt)
	isExhausted := otp.Attempts >= otp.MaxAttempts
	if otp.IsUsed {
		return "", ErrUsedOTPRequest
	} else if isExpired {
		return "", ErrExpiredOTPRequest
	} else if isExhausted {
		return "", ErrReachMaxAttempt
	}

	return otp.Token, nil
}

func (r *OTPRepository) IncrOTPRequestAttempt(ctx context.Context, tel string, ref string) error {
	result := r.db.WithContext(ctx).
		Model(&models.OtpRequest{}).
		Where("tel = ? AND ref_no = ? AND is_used = ?", tel, ref, false).
		Update("attempts", gorm.Expr("attempts + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OTPRepository) HasRecentRequestByTel(ctx context.Context, tel string) (bool, error) {
	var count int64
	threshold := time.Now().Add(-2 * time.Minute)

	err := r.db.WithContext(ctx).
		Model(&models.OtpRequest{}).
		Where("tel = ? AND created_at > ?", tel, threshold).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *OTPRepository) IsUserExceedRequestLimit(ctx context.Context, userID uint) (bool, error) {
	var count int64
	threshold := time.Now().Add(-UserRequestWindowPeriod)

	err := r.db.WithContext(ctx).
		Model(&models.OtpRequest{}).
		Where("user_id = ? AND created_at > ?", userID, threshold).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count >= MaxRequestByUser, nil
}

func (r *OTPRepository) MarkOTPRequestUsed(ctx context.Context, tel string, ref string) (uint, error) {
	var otp models.OtpRequest

	if err := r.db.WithContext(ctx).
		Where("tel = ? AND ref_no = ? AND is_used = false AND expires_at > ? AND attempts < max_attempts", tel, ref, time.Now()).
		First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrExpiredOTPRequest
		}
		return 0, err
	}

	if err := r.db.WithContext(ctx).
		Model(&otp).
		Update("is_used", true).Error; err != nil {
		return 0, err
	}

	return otp.UserID, nil
}
