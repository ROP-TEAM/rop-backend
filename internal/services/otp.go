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
	"regexp"
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

const (
	StatusSuccess = "success"
)

type RequestOTPRequest struct {
	Tel    string `json:"tel" example:"0812345678"`
	UserId uint   `json:"user_id" example:"1"`
}

type RequestOTPResponse struct {
	RefNo  string `json:"refNo" example:"ABC123"`
	Status string `json:"status" example:"success"`
}

type requestThaiBulkSuccessResponse struct {
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
	Pin   string `json:"pin" example:"123456"`
	Tel   string `json:"tel" example:"0812345678"`
	RefNo string `json:"refNo" example:"ABC123"`
}

type VerifyOTPResponse struct {
	Status string `json:"status" example:"success"`
}

type OTPService struct {
	userRepository    *repository.UserRepository
	otpRepository     *repository.OTPRepository
	cfg               *config.Config
	httpClient        *http.Client
	db                *gorm.DB
	mockOtpRepository *repository.MockOTPRepository
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
	mockOtpRepository := repository.NewMockOTPRepository(db)
	return &OTPService{
		userRepository: userRepository,
		otpRepository:  otpRepository,
		cfg:            cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		db:                db,
		mockOtpRepository: mockOtpRepository,
	}
}

func (s *OTPService) RequestOTP(ctx context.Context, userID uint, req *RequestOTPRequest) (*RequestOTPResponse, error) {
	user, err := s.userRepository.FindByPhone(ctx, req.Tel)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("OTPRequest: database cannot find %v: %v", req.Tel, err)
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

	// apiRes, err := s.callThaiBulkRequestOTP(ctx, req.Tel)

	// mock for development

	apiRes, err := s.callMockThaiBulkRequestOTP(ctx, req.Tel)

	if err != nil {
		log.Printf("OTPRequest: fetching thaibulk api : %v", err)
		return nil, err
	}

	//at this point success on calling thaibulk api

	if apiRes.Status != StatusSuccess {
		return &RequestOTPResponse{RefNo: "", Status: apiRes.Status}, nil
	}

	err = s.otpRepository.CreateOTPRequest(ctx, &repository.CreateOTPRequest{
		Tel:    req.Tel,
		Token:  apiRes.Token,
		RefNo:  apiRes.RefNo,
		UserId: userID,
	})
	if err != nil {
		log.Printf("OTPRequest:  database saving otp request: %v", err)
		return nil, err
	}

	return &RequestOTPResponse{RefNo: apiRes.RefNo, Status: apiRes.Status}, nil
}

func (s *OTPService) callThaiBulkRequestOTP(ctx context.Context, tel string) (*requestThaiBulkSuccessResponse, error) {
	vals := url.Values{}
	vals.Set("key", s.cfg.OTP_APP_KEY)
	vals.Set("secret", s.cfg.OTP_APP_SECRET)
	vals.Set("msisdn", tel)

	payload := strings.NewReader(vals.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.OTP_APP_URL_REQUEST, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request body: %w", err)
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

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		var apiSuccess requestThaiBulkSuccessResponse
		if err := json.Unmarshal(body, &apiSuccess); err != nil {
			return nil, fmt.Errorf("parsing success response: %w", err)
		}
		return &apiSuccess, nil
	}

	//error response
	var apiError thaiBulkErrorResponse
	if err := json.Unmarshal(body, &apiError); err != nil {
		return nil, fmt.Errorf("%w: invalid error response (status=%d): %v", ErrOTPProvider, res.StatusCode, err)
	}

	errMsg := "unknown error"
	if len(apiError.Errors) > 0 {
		errMsg = apiError.Errors[0].Message
	}
	return nil, fmt.Errorf("%w: %s (status=%d)", ErrOTPProvider, errMsg, res.StatusCode)
}

func (s *OTPService) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error) {

	token, err := s.otpRepository.IsOTPRequestMutable(ctx, req.Tel, req.RefNo)
	if err != nil {
		log.Printf("OTPVerify: database pre-checkck: %v", err)
		return nil, err
	}

	if token == "" {
		log.Printf("OTPVerify: empty token retrieved for tel=%s ref=%s", req.Tel, req.RefNo)
		return nil, fmt.Errorf("invalid OTP token format")
	}

	// apiRes, err := s.callThaiBulkVerifyAPI(ctx, token, req.Pin)

	// for testing

	apiRes, err := s.callMockThaiBulkVerifyOTP(ctx, token, req.Pin)

	if err != nil {
		if incrErr := s.otpRepository.IncrOTPRequestAttempt(ctx, req.Tel, req.RefNo); incrErr != nil {
			log.Printf("OTPVerify: failed to increment attempt: %v", incrErr)
		}
		return nil, err
	}

	if apiRes.Status != StatusSuccess {
		if err := s.otpRepository.IncrOTPRequestAttempt(ctx, req.Tel, req.RefNo); err != nil {
			log.Printf("OTPVerify: error when increase attempt after enter invalid pin: %v", err)
			return nil, err
		}
		return nil, ErrInvalidOTP
	}

	// success
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userID, err := s.otpRepository.WithTx(tx).MarkOTPRequestUsed(ctx, req.Tel, req.RefNo)
		if err != nil {
			log.Printf("OTPVerify: failed to mark OTP as used: %v", err)
			return fmt.Errorf("marking OTP as used: %w", err)
		}

		if err := s.userRepository.WithTx(tx).CompleteUserValidation(ctx, int(userID), req.Tel); err != nil {
			log.Printf("OTPVerify: failed to complete user validation: %v", err)
			return fmt.Errorf("completing user validation: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("OTPVerify: transaction failed: %v", err)
		return nil, fmt.Errorf("completing verification: %w", err)
	}

	return &VerifyOTPResponse{Status: StatusSuccess}, nil

}

func (s *OTPService) callThaiBulkVerifyAPI(ctx context.Context, token, pin string) (*verifyThaiBulkSuccessResponse, error) {
	vals := url.Values{}
	vals.Set("key", s.cfg.OTP_APP_KEY)
	vals.Set("secret", s.cfg.OTP_APP_SECRET)
	vals.Set("token", token)
	vals.Set("pin", pin)

	payload := strings.NewReader(vals.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.OTP_APP_URL_VERIFY, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request body: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %v", ErrOTPProvider, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	//at this point success on calling thaibulk api

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		var apiSuccess verifyThaiBulkSuccessResponse
		if err := json.Unmarshal(body, &apiSuccess); err != nil {
			return nil, fmt.Errorf("parsing success response: %w", err)
		}
		return &apiSuccess, nil
	}

	//error response
	var apiError thaiBulkErrorResponse
	if err := json.Unmarshal(body, &apiError); err != nil {
		return nil, fmt.Errorf("%w: invalid error response (status=%d): %v", ErrOTPProvider, res.StatusCode, err)
	}

	errMsg := "unknown error"
	if len(apiError.Errors) > 0 {
		errMsg = apiError.Errors[0].Message
	}
	return nil, fmt.Errorf("%w: %s (status=%d)", ErrOTPProvider, errMsg, res.StatusCode)
}

//
// mock for testing and development
//

var thaiTelRegex = regexp.MustCompile(`^(06|08|09)[0-9]{8}$`)

func isValidThaiTel(tel string) bool {
	return thaiTelRegex.MatchString(tel)
}

func (s *OTPService) callMockThaiBulkRequestOTP(ctx context.Context, tel string) (*requestThaiBulkSuccessResponse, error) {

	if !isValidThaiTel(tel) {
		return nil, fmt.Errorf("%w: %s (status=%d)", ErrOTPProvider, "Gateway response send sms fail.", 400)
	}

	otp, err := s.mockOtpRepository.CreateOTP(ctx, tel)
	if err != nil {
		return nil, fmt.Errorf("%w: %s (status=%d)", ErrOTPProvider, "Gateway response send sms fail.", 400)
	}

	log.Printf("\nSOROUTETION OTP: %s\n(Ref: %s) Valid for %d mintutes.\n", otp.Pin, otp.RefNo, repository.OTPExpiredTime)
	return &requestThaiBulkSuccessResponse{
		Status: StatusSuccess,
		Token:  otp.Token,
		RefNo:  otp.RefNo,
	}, nil
}

func (s *OTPService) callMockThaiBulkVerifyOTP(
	ctx context.Context,
	token string,
	pin string,
) (*verifyThaiBulkSuccessResponse, error) {

	result, reason, err := s.mockOtpRepository.VerifyOTP(ctx, token, pin)
	if err != nil {
		return nil, err // real system error
	}

	if !result {
		return &verifyThaiBulkSuccessResponse{
			Status:  "fail",
			Message: string(reason),
		}, nil
	}

	return &verifyThaiBulkSuccessResponse{
		Status:  "success",
		Message: "OTP verified successfully",
	}, nil
}
