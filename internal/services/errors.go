package services

//domain error for service to handler
import (
	"errors"
	"time"
)

var (
	ErrUsedPhoneNumber     = errors.New("phone number already in use")
	ErrOTPProvider         = errors.New("otp provider error")
	ErrUserReachMaxRequest = errors.New("user reach max request")
	ErrInvalidOTP          = errors.New("invalid otp pin")
	ErrOTPNotFound         = errors.New("there is no OTP request for this user and tel")
	ErrReachMaxAttempt     = errors.New("verification reach the maximum attempt")
	ErrExpiredOTPRequest   = errors.New("this otpRequest already expired")
	ErrUsedOTPRequest      = errors.New("tel already be verified by otp pin")
)

type ErrPhoneNumberHasRecentRequest struct {
	RetryAfter time.Time
}

func (e *ErrPhoneNumberHasRecentRequest) Error() string {
	return "phone number has recent request"
}
