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
var ErrPhoneNumberHasRecentRequest = errors.New("phone number has recent request")
var ErrInvalidOTP = errors.New("invalid otp pin")

type RequestOTPRequest struct {
	Tel    string `json:"tel"`
	UserId uint   `json:"user_id"`
}

type RequestOTPResponse struct {
	RefNo string `json:"refNo"`
}

type thaiBulkSuccessResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	RefNo  string `json:"refno"`
}

type verifyThaiBulkSuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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

type VerifyOTPRequest struct {
	Pin   string `json:"pin"`
	Tel   string `json:"tel"`
	RefNo string `json:"refNo"`
}

type VerifyOTPResponse struct {
	Status string `json:"status"`
}

type OTPService struct {
	userRepository *repository.UserRepository
	otpRepository  *repository.OTPRepository
	cfg            *config.Config
	httpClient     *http.Client
	db             *gorm.DB
}

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

func NewOTPService(db *gorm.DB, cfg *config.Config) *OTPService {
	userRepository := repository.NewUserRepository(db)
	otpRepository := repository.NewOTPRepository(db)
	return &OTPService{
		userRepository: userRepository,
		otpRepository:  otpRepository,
		cfg:            cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		db: db,
	}
}

func (s *OTPService) RequestOTP(ctx context.Context, req *RequestOTPRequest) (*RequestOTPResponse, error) {
	user, err := s.userRepository.FindByPhone(ctx, req.Tel)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("OTPRequest: database error: %v", err)
		return nil, err
	}

	if user != nil && user.IsValidated {
		log.Printf("OTPRequest: database bad requests, phone already use: %v", req.Tel)
		return nil, ErrUsedPhoneNumber
	}

	hasRecentRequest, err := s.otpRepository.HasRecentRequest(ctx, req.Tel)
	if err != nil {
		log.Printf("OTPRequest: database error: %v", err)
		return nil, err
	}
	if hasRecentRequest {
		log.Printf("OTPRequest: phone has recent request : %v", req.Tel)
		return nil, ErrPhoneNumberHasRecentRequest
	}

	vals := url.Values{}
	vals.Set("key", s.cfg.OTP_APP_KEY)
	vals.Set("secret", s.cfg.OTP_APP_SECRET)
	vals.Set("msisdn", req.Tel)
	payload := strings.NewReader(vals.Encode())
	thaibulkReq, err := http.NewRequestWithContext(ctx, "POST", s.cfg.OTP_APP_URL_REQUEST, payload)
	if err != nil {
		log.Printf("OTPRequest: error when create new request: %v", err)
		return nil, err
	}

	thaibulkReq.Header.Set("accept", "application/json")
	thaibulkReq.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := s.httpClient.Do(thaibulkReq)
	if err != nil {
		log.Printf("OTPRequest: provider request failed: tel=%s err=%v", req.Tel, err)
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

		err := s.otpRepository.CreateOTPRequest(ctx, &repository.CreateOTPRequest{
			Tel:    req.Tel,
			Token:  apiSuccess.Token,
			RefNo:  apiSuccess.RefNo,
			UserId: req.UserId,
		})
		if err != nil {
			log.Printf("OTPRequest: error when save otp request: %v", err)
			return nil, err
		}

		return &RequestOTPResponse{RefNo: apiSuccess.RefNo}, nil
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

func (s *OTPService) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error) {

	token, err := s.otpRepository.IsOTPRequestMutable(ctx, req.Tel, req.RefNo)
	if err != nil {
		log.Printf("OTPVerify: database pre-checkck: %v", err)
		return nil, err
	}

	vals := url.Values{}
	vals.Set("key", s.cfg.OTP_APP_KEY)
	vals.Set("secret", s.cfg.OTP_APP_SECRET)
	vals.Set("token", token)
	vals.Set("pin", req.Pin)

	payload := strings.NewReader(vals.Encode())

	thaibulkReq, err := http.NewRequestWithContext(ctx, "POST", s.cfg.OTP_APP_URL_VERIFY, payload)
	if err != nil {
		log.Printf("OTPVerify: error when create new verify: %v", err)
		return nil, err
	}

	thaibulkReq.Header.Set("accept", "application/json")
	thaibulkReq.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := s.httpClient.Do(thaibulkReq)
	if err != nil {
		log.Printf("OTPVerify: provider request failed: tel=%s err=%v", req.Tel, err)
		return nil, fmt.Errorf("%w: request failed: %v", ErrOTPProvider, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("OTPVerify: error when read response body: %v", err)
		return nil, err
	}

	//at this point success on calling thaibulk api

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		var apiSuccess verifyThaiBulkSuccessResponse
		if err := json.Unmarshal(body, &apiSuccess); err != nil {
			log.Printf("OTPVerify: error when read unmarshal body: %v", err)
			return nil, err
		}

		if apiSuccess.Status != "success" {
			if err := s.otpRepository.IncrOTPRequestAttempt(ctx, req.Tel, req.RefNo); err != nil {
				log.Printf("OTPVerify: error when increase attempt: %v", err)
				return nil, err
			}
			return nil, ErrInvalidOTP
		}

		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			userID, err := s.otpRepository.WithTx(tx).MarkOTPRequestUsed(ctx, req.Tel, req.RefNo)
			if err != nil {
				log.Printf("OTPVerify: error when save otp request: %v", err)

				return err
			}
			return s.userRepository.WithTx(tx).CompleteUserValidation(ctx, int(userID), req.Tel)
		})
		if err != nil {
			log.Printf("OTPVerify: error in commit transaction: %v", err)
			return nil, err
		}

		return &VerifyOTPResponse{Status: apiSuccess.Status}, nil
	}

	//fail code
	err = s.otpRepository.IncrOTPRequestAttempt(ctx, req.Tel, req.RefNo)
	if err != nil {
		log.Printf("OTPVerify: error when increase attempt: %v", err)
		return nil, err
	}

	var apiError thaiBulkErrorResponse
	if err := json.Unmarshal(body, &apiError); err != nil {
		log.Printf("OTPVerify: error when read unmarshal body: %v", err)
		return nil, fmt.Errorf("%w: invalid error response: %v", ErrOTPProvider, err)
	}

	errMsg := "unknown error"
	if len(apiError.Errors) > 0 {
		errMsg = apiError.Errors[0].Message
	}

	log.Printf("OTPVerify: error from api status fail: %v", errMsg)
	return nil, fmt.Errorf("%w: %s (status=%d)", ErrOTPProvider, errMsg, res.StatusCode)

}
