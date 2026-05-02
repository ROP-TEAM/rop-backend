package handlers

import (
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/services"

	"context"
	"errors"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

var phoneRegex = regexp.MustCompile(`^\d{10}$`)

func respondSuccess(c fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(OTPResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func respondError(c fiber.Ctx, status int, message string, code string, data interface{}) error {
	return c.Status(status).JSON(OTPResponse{
		Success: false,
		Message: message,
		Data:    data,
		Error:   &OTPErrorInfo{Code: code},
	})
}

type OTPHandler struct {
	service *services.OTPService
}

func NewOTPHandler(service *services.OTPService) *OTPHandler {
	return &OTPHandler{service: service}
}

// RequestOTP godoc
// @Summary      Request OTP SMS
// @Description  Send OTP to a phone number and receive token + refno
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      services.RequestOTPRequest  true  "Phone number payload"
// @Success      200      {object}  handlers.APIResponse{data=services.RequestOTPResponse}
// @Failure      400      {object}  handlers.APIResponse
// @Failure      409      {object}  handlers.APIResponse
// @Failure      502      {object}  handlers.APIResponse
// @Failure      500      {object}  handlers.APIResponse
// @Router       /api/auth/otp [post]
func (h *OTPHandler) RequestOTP(c fiber.Ctx) error {
	var body models.RequestOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body", "BAD_REQUEST", nil)
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return respondError(c, fiber.StatusUnauthorized, "user unauthorized", "BAD_REQUEST", nil)
	}

	userID := claims.UserID

	if body.Tel == "" {
		return respondError(c, fiber.StatusBadRequest, "tel field is required", "BAD_REQUEST", nil)
	}

	if !phoneRegex.MatchString(body.Tel) {
		return respondError(c, fiber.StatusBadRequest, "invalid tel format", "BAD_REQUEST", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.service.RequestOTP(ctx, userID, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUsedPhoneNumber):
			return respondError(c, fiber.StatusConflict, "tel already in used", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrPhoneNumberHasRecentRequest):
			return respondError(c, fiber.StatusConflict, "phone has recently request for otp", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrUserReachMaxRequest):
			return respondError(c, fiber.StatusConflict, "user reach max request for OTP in the period", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrOTPProvider):
			return respondError(c, fiber.StatusBadGateway, "OTP service unavailable", "OTP_UNAVAILABLE", nil)

		default:
			return respondError(c, fiber.StatusInternalServerError, "internal server error", "INTERNAL_ERROR", nil)
		}
	}

	return respondSuccess(c, fiber.StatusOK, "successful request for otp", res)
}

// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verify OTP pin using refNo and phone number. Handles expiration, attempt limits, and reuse protection.
// @Tags         auth
// @ID           verify-otp
// @Accept       json
// @Produce      json
// @Param        request  body      services.VerifyOTPRequest  true  "OTP verification payload"
// @Success      200  {object}  handlers.APIResponse{data=services.VerifyOTPResponse} "OTP verified successfully"
// @Failure      400  {object}  handlers.APIResponse "Invalid request body or missing fields"
// @Failure      401  {object}  handlers.APIResponse "Invalid OTP PIN"
// @Failure      409  {object}  handlers.APIResponse "Expired / already used / max attempts reached"
// @Failure      500  {object}  handlers.APIResponse "Internal server error"
// @Router       /api/auth/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c fiber.Ctx) error {
	var body models.VerifyOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body", "BAD_REQUEST", nil)
	}

	if body.Tel == "" {
		return respondError(c, fiber.StatusBadRequest, "tel field is required", "BAD_REQUEST", nil)
	}

	if !phoneRegex.MatchString(body.Tel) {
		return respondError(c, fiber.StatusBadRequest, "invalid tel format", "BAD_REQUEST", nil)
	}

	if body.Pin == "" {
		return respondError(c, fiber.StatusBadRequest, "pin field is required", "BAD_REQUEST", nil)
	}

	if body.RefNo == "" {
		return respondError(c, fiber.StatusBadRequest, "refNo field is required", "BAD_REQUEST", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.service.VerifyOTP(ctx, &body)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReachMaxAttempt):
			return respondError(c, fiber.StatusConflict, "number of attempt already reach the limit", "OTP_LIMITED", nil)

		case errors.Is(err, repository.ErrExpiredOTPRequest):
			return respondError(c, fiber.StatusConflict, "OTP already expire for verification", "OTP_LIMITED", nil)

		case errors.Is(err, repository.ErrUsedOTPRequest):
			return respondError(c, fiber.StatusConflict, "OTP verification already done", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrInvalidOTP):
			return respondError(c, fiber.StatusUnauthorized, "invalid OTP pin", "OTP_LIMITED", nil)

		case errors.Is(err, gorm.ErrRecordNotFound):
			return respondError(c, fiber.StatusUnauthorized, "record according to tel and refNo not found", "BAD_REQUEST", nil)

		default:
			return respondError(c, fiber.StatusInternalServerError, "internal server error", "INTERNAL_ERROR", nil)
		}
	}
	return respondSuccess(c, fiber.StatusOK, "OTP verified successfully", res)
}
