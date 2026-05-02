package handlers

import (
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/services"

	"context"
	"errors"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
)

var phoneRegex = regexp.MustCompile(`^\d{10}$`)

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
	var body services.RequestOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("invalid request body"))
	}

	claims := middleware.GetUser(c)
	if claims == nil {
		return c.Status(401).JSON(Error("unauthorized"))
	}

	userID := claims.UserID

	if body.Tel == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("tel field is required"))
	}

	// if body.UserId == "" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
	// 		Error: "user_id field is required",
	// 	})
	// }

	if !phoneRegex.MatchString(body.Tel) {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("invalid tel format"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.service.RequestOTP(ctx, userID, &body)
	if err != nil {

		if errors.Is(err, services.ErrUsedPhoneNumber) {
			return c.Status(fiber.StatusConflict).JSON(Fail("tel already in use"))
		}

		if errors.Is(err, services.ErrPhoneNumberHasRecentRequest) {
			return c.Status(fiber.StatusConflict).JSON(Fail("tel recently request for OTP"))
		}

		if errors.Is(err, services.ErrOTPProvider) {
			return c.Status(fiber.StatusBadGateway).JSON(Fail("OTP service unavailable"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Error("internal server error"))
	}

	return c.JSON(Success("", res))
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
	var body services.VerifyOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Error("invalid request body"))
	}

	if body.Tel == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("tel field is required"))
	}

	if !phoneRegex.MatchString(body.Tel) {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("invalid tel format"))
	}

	if body.Pin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("Pin field is required"))
	}

	if body.RefNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Fail("refNo field is required"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.service.VerifyOTP(ctx, &body)
	if err != nil {
		if errors.Is(err, repository.ErrReachMaxAttempt) {
			return c.Status(fiber.StatusConflict).JSON(Fail("number of attempt already reach the limit"))
		}

		if errors.Is(err, repository.ErrExpiredOTPRequest) {
			return c.Status(fiber.StatusConflict).JSON(Fail("OTP already expire for verification"))
		}

		if errors.Is(err, repository.ErrUsedOTPRequest) {
			return c.Status(fiber.StatusConflict).JSON(Fail("OTP verification already done"))
		}

		if errors.Is(err, services.ErrInvalidOTP) {
			return c.Status(fiber.StatusUnauthorized).JSON(Fail("invalid OTP pin"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Error("internal server error"))
	}
	return c.JSON(Success("", res))
}
