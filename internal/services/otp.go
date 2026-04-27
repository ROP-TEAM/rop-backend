package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"strings"

	"ROP_Backend/internal/config"
	"ROP_Backend/internal/repository"

	"gorm.io/gorm"
)

var ErrUsedPhoneNumber = errors.New("phone number already in use")
var ErrOTPProvider = errors.New("otp provider error")

type RequestOTPResponse struct {
	Token string `json:"token"`
	RefNo string `json:"refno"`
}

type thaiBulkSuccessResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	RefNo  string `json:"refno"`
}

type thaiBulkError struct {
	Detail  StringOrArray `json:"detail"`
	Message string        `json:"message"`
}

type thaiBulkErrorResponse struct {
	Errors []thaiBulkError `json:"errors"`
	Code   int             `json:"code"`
}

type StringOrArray []string

// helper function for ThaiBulkErrorResponse
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	// try array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*s = arr
		return nil
	}

	// try string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = []string{str}
		return nil
	}

	return fmt.Errorf("invalid detail format")
}

type OTPService struct {
	userRepository *repository.UserRepository
	cfg            *config.Config
	httpClient     *http.Client
}

func NewOTPService(db *gorm.DB, cfg *config.Config) *OTPService {
	userRepository := repository.NewUserRepository(db)
	return &OTPService{
		userRepository: userRepository,
		cfg:            cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *OTPService) RequestOTP(ctx context.Context, tel string) (*RequestOTPResponse, error) {
	user, err := s.userRepository.FindByPhone(ctx, tel)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("OTPRequest: database error: %v", err)
		return nil, err
	}

	if user != nil && user.IsValidated {
		log.Printf("OTPRequest: database bad requests, phone already use: %v", tel)
		return nil, ErrUsedPhoneNumber
	}

	vals := url.Values{}
	vals.Set("key", s.cfg.OTP_APP_KEY)
	vals.Set("secret", s.cfg.OTP_APP_SECRET)
	vals.Set("msisdn", tel)
	payload := strings.NewReader(vals.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.OTP_APP_URL, payload)
	if err != nil {
		log.Printf("OTPRequest: error when create new request: %v", err)
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("OTPRequest: provider request failed: tel=%s err=%v", tel, err)
		return nil, fmt.Errorf("%w: request failed: %v", ErrOTPProvider, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("OTPRequest: error when read response body: %v", err)
		return nil, err
	}

	//at this point success on calling thaibulk api
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		var apiSuccess thaiBulkSuccessResponse
		if err := json.Unmarshal(body, &apiSuccess); err != nil {
			log.Printf("OTPRequest: error when read unmarshal body: %v", err)
			return nil, err
		}

		if apiSuccess.Status != "success" {
			return nil, fmt.Errorf("%w: unexpected status %q", ErrOTPProvider, apiSuccess.Status)
		}

		return &RequestOTPResponse{Token: apiSuccess.Token, RefNo: apiSuccess.RefNo}, nil
	}

	//fail code
	var apiError thaiBulkErrorResponse
	if err := json.Unmarshal(body, &apiError); err != nil {
		log.Printf("OTPRequest: error when read unmarshal body: %v", err)
		return nil, fmt.Errorf("%w: invalid error response: %v", ErrOTPProvider, err)
	}

	var errMsg string
	if len(apiError.Errors) > 0 {
		errMsg = apiError.Errors[0].Message
	} else {
		errMsg = "unknown error"
	}
	log.Printf("OTPRequest: error from api status fail: %v", errMsg)
	return nil, fmt.Errorf("%w: %s (status=%d)", ErrOTPProvider, errMsg, res.StatusCode)
}

type VerifyOTPRequest struct {
	OTP   string `json:"otp"`
	Token string `json:"token"`
}

type verifyOTPResponse struct {
}

func (s *OTPService) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*verifyOTPResponse, error) {
	return nil, nil
}
