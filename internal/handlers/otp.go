package handlers

import (
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/services"
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
)

type OTPHandler struct {
	service *services.OTPService
}

func NewOTPHandler(service *services.OTPService) *OTPHandler {
	return &OTPHandler{service: service}
}

var phoneRegex = regexp.MustCompile(`^\d{10}$`)

// RequestOTP godoc
// @Summary      Request OTP SMS
// @Description  Sends an OTP to the provided phone number and returns a token and reference number for verification.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      handlers.otpRequest  true  "Phone number in 10-digit format"
// @Success      200      {object}  services.OTPResponse "Successfully sent OTP"
// @Failure      400      {object}  handlers.ErrorResponse "Invalid phone format or request body"
// @Failure      409      {object}  handlers.ErrorResponse "Phone number already validated"
// @Failure      502      {object}  handlers.ErrorResponse "OTP Provider (ThaiBulk) error"
// @Failure      500      {object}  handlers.ErrorResponse "Internal server error"
// @Router       /api/auth/otp [post]
func (h *OTPHandler) RequestOTP(c fiber.Ctx) error {
	var body services.RequestOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "invalid request body",
		})
	}

	if body.Tel == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "tel field is required",
		})
	}

	// if body.UserId == "" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
	// 		Error: "user_id field is required",
	// 	})
	// }

	if !phoneRegex.MatchString(body.Tel) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "invalid tel format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.service.RequestOTP(ctx, &body)
	if err != nil {

		if errors.Is(err, services.ErrUsedPhoneNumber) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "phone number already in use",
			})
		}

		if errors.Is(err, services.ErrPhoneNumberHasRecentRequest) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "phone number recently request for OTP",
			})
		}

		if errors.Is(err, services.ErrOTPProvider) {
			return c.Status(fiber.StatusBadGateway).JSON(ErrorResponse{
				Error: "OTP service unavailable",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.JSON(res)
}

func (h *OTPHandler) VerifyOTP(c fiber.Ctx) error {
	var body services.VerifyOTPRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "invalid request body",
		})
	}

	if body.Tel == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "tel field is required",
		})
	}

	if !phoneRegex.MatchString(body.Tel) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "invalid tel format",
		})
	}

	if body.Pin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Pin field is required",
		})
	}

	if body.RefNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "refNo field is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := h.service.VerifyOTP(ctx, &body)
	if err != nil {
		if errors.Is(err, repository.ErrReachMaxAttempt) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "number of attempt already reach the limit",
			})
		}

		if errors.Is(err, repository.ErrExpiredOTPRequest) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "OTP already expire for verification",
			})
		}

		if errors.Is(err, repository.ErrUsedOTPRequest) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "OTP verification already done",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "internal server error",
		})
	}
	return c.JSON(res)
}
