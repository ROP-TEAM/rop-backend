package services

//domain error for service to handler
import (
	"errors"
	"time"
)

var (
	ErrUsedPhoneNumber      = errors.New("phone number already in use")
	ErrOTPProvider          = errors.New("otp provider error")
	ErrUserReachMaxRequest  = errors.New("user reach max request")
	ErrInvalidOTP           = errors.New("invalid otp pin")
	ErrOTPNotFound          = errors.New("there is no OTP request for this user and tel")
	ErrReachMaxAttempt      = errors.New("verification reach the maximum attempt")
	ErrExpiredOTPRequest    = errors.New("this otpRequest already expired")
	ErrUsedOTPRequest       = errors.New("tel already be verified by otp pin")
	ErrUserNotFound         = errors.New("user not found")
	ErrPlanNotFound         = errors.New("plan not found")
	ErrUserHasNoCompany     = errors.New("user has not been employed yet")
	ErrCreatingPlan         = errors.New("creating plan")
	ErrCreatingTagSkill     = errors.New("creating tag skill")
	ErrCreatingOrder        = errors.New("creating order")
	ErrCreatingVehicle      = errors.New("creating vehicle")
	ErrCreatingOrderSkill   = errors.New("creating order tag skill")
	ErrCreatingVehicleSkill = errors.New("creating vehicl tag skill")
	ErrInvalidUser          = errors.New("invalid user or user not found ")
	ErrInvalidDateFormat    = errors.New("invalid date format ")
	ErrPlanNameTooLong      = errors.New("plan name must not exceed 100 characters")
)

type ErrPhoneNumberHasRecentRequest struct {
	RetryAfter time.Time
}

func (e *ErrPhoneNumberHasRecentRequest) Error() string {
	return "phone number has recent request"
}
