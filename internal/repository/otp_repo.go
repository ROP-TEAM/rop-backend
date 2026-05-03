package repository

import (
	"ROP_Backend/internal/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNoRowsAffected = errors.New("no rows affected")
)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) WithTx(tx *gorm.DB) *OTPRepository {
	return &OTPRepository{db: tx}
}

func (r *OTPRepository) Create(ctx context.Context, req *models.OtpRequest) (*models.OtpRequest, error) {
	err := r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *OTPRepository) FindByTelAndRef(ctx context.Context, tel string, ref string) (*models.OtpRequest, error) {
	var otp models.OtpRequest
	err := r.db.WithContext(ctx).
		Where("tel = ? AND ref_no = ?", tel, ref).
		Order("created_at DESC").
		First(&otp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &otp, err
}

func (r *OTPRepository) IncrAttemptByID(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.OtpRequest{}).
		Where("id = ?", id).
		Update("attempts", gorm.Expr("attempts + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *OTPRepository) FindLatestByTelSince(ctx context.Context, tel string, since time.Time) (*models.OtpRequest, error) {
	var otp models.OtpRequest
	err := r.db.WithContext(ctx).
		Where("tel = ? AND created_at > ?", tel, since).
		Order("created_at DESC").
		First(&otp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &otp, err
}

func (r *OTPRepository) CountByUserIDSince(ctx context.Context, userID uint, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.OtpRequest{}).
		Where("user_id = ? AND created_at > ?", userID, since).
		Count(&count).Error
	return count, err
}

func (r *OTPRepository) MarkUsedByID(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.OtpRequest{}).
		Where("id = ? AND is_used = false AND expires_at > ? AND attempts < max_attempts", id, time.Now()).
		Update("is_used", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}
