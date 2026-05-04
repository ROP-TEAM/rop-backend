package handlers

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"
	"ROP_Backend/internal/utils"

	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
)

var (
	recentRequestErr *services.ErrPhoneNumberHasRecentRequest
)

func respondSuccess(c fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(dto.OTPResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func respondError(c fiber.Ctx, status int, message string, code string, data interface{}) error {
	return c.Status(status).JSON(dto.OTPResponse{
		Success: false,
		Message: message,
		Data:    data,
		Error:   &dto.OTPErrorInfo{Code: code},
	})
}

type OTPHandler struct {
	service *services.OTPService
}

func NewOTPHandler(service *services.OTPService) *OTPHandler {
	return &OTPHandler{service: service}
}

// RequestOTP godoc
// @Summary      Request OTP
// @Description  Send an OTP SMS to the given phone number. Requires a valid JWT. Returns a refNo to use during verification.
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RequestOTPRequest  true  "Phone number payload"
// @Success      200      {object}  dto.OTPResponse{data=dto.RequestOTPResponse}  "OTP sent successfully"
// @Failure      400      {object}  dto.OTPResponse  "Missing/invalid tel field or malformed body"
// @Failure      401      {object}  dto.OTPResponse  "Missing or invalid JWT"
// @Failure      409      {object}  dto.OTPResponse  "Phone already verified / recent request exists / user exceeded request limit"
// @Failure      502      {object}  dto.OTPResponse  "Upstream OTP provider unavailable"
// @Failure      500      {object}  dto.OTPResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /api/onboarding/otp [post]
func (h *OTPHandler) RequestOTP(c fiber.Ctx) error {
	var body dto.RequestOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body", "BAD_REQUEST", map[string]any{"details": err})
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return respondError(c, fiber.StatusUnauthorized, "user unauthorized", "BAD_REQUEST", nil)
	}

	userID := claims.UserID

	if body.Tel == "" {
		return respondError(c, fiber.StatusBadRequest, "tel field is required", "BAD_REQUEST", nil)
	}

	if !utils.IsThaiMobile(body.Tel) {
		return respondError(c, fiber.StatusBadRequest, "invalid tel format", "BAD_REQUEST", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.service.RequestOTP(ctx, userID, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUsedPhoneNumber):
			return respondError(c, fiber.StatusConflict, "tel already in used", "OTP_LIMITED", nil)

		case errors.As(err, &recentRequestErr):
			return respondError(c, fiber.StatusConflict, "phone has recently request for otp", "OTP_LIMITED", map[string]any{"retry_after": recentRequestErr.RetryAfter})

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
// @Description  Verify an OTP pin using the refNo and phone number from the request step. Increments attempt count on failure and marks the OTP as used on success.
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request  body      dto.VerifyOTPRequest  true  "OTP verification payload"
// @Success      200      {object}  dto.OTPResponse{data=dto.VerifyOTPResponse}  "OTP verified; user phone validated"
// @Failure      400      {object}  dto.OTPResponse  "Missing/invalid tel, pin, or refNo — or record not found"
// @Failure      401      {object}  dto.OTPResponse  "Incorrect OTP pin"
// @Failure      409      {object}  dto.OTPResponse  "OTP expired / already used / max attempts reached"
// @Failure      500      {object}  dto.OTPResponse  "Internal server error"
// @Router       /api/onboarding/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c fiber.Ctx) error {
	var body dto.VerifyOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body", "BAD_REQUEST", nil)
	}

	if body.Tel == "" {
		return respondError(c, fiber.StatusBadRequest, "tel field is required", "BAD_REQUEST", nil)
	}

	if !utils.IsThaiMobile(body.Tel) {
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
		case errors.Is(err, services.ErrReachMaxAttempt):
			return respondError(c, fiber.StatusConflict, "number of attempt already reach the limit", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrExpiredOTPRequest):
			return respondError(c, fiber.StatusConflict, "OTP already expire for verification", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrUsedOTPRequest):
			return respondError(c, fiber.StatusConflict, "OTP verification already done", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrInvalidOTP):
			return respondError(c, fiber.StatusUnauthorized, "invalid OTP pin", "OTP_LIMITED", nil)

		case errors.Is(err, services.ErrOTPNotFound):
			return respondError(c, fiber.StatusUnauthorized, "record according to tel and refNo not found", "BAD_REQUEST", nil)

		default:
			return respondError(c, fiber.StatusInternalServerError, "internal server error", "INTERNAL_ERROR", nil)
		}
	}
	return respondSuccess(c, fiber.StatusOK, "OTP verified successfully", res)
}
