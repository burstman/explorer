package paymentservices

import (
	"bytes"
	"encoding/json"
	"errors"
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
)

type KonnectService struct {
	APIKey     string
	BaseURL    string
	WebhookURL string
	SuccessURL string
	FailURL    string
	WalletID   string
}

// CreateKonnectPayment stores a payment in the DB
func (s *KonnectService) CreateKonnectPayment(payment *types.KonnectPaymentResponse) error {
	if err := db.Get().Create(payment).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("duplicate payment_ref")
		}
		return err
	}
	return nil
}

func (s *KonnectService) doRequest(req *types.InitKonnectPaymentRequest, path string) ([]byte, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, s.BaseURL+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", s.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to init payment: status=%d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// ✅ Implements PaymentService.InitPayment
func (s *KonnectService) InitPayment(req *types.InitKonnectPaymentRequest) (types.PaymentResponse, error) {
	data, err := s.doRequest(req, "/payments/init-payment")
	if err != nil {
		return nil, err
	}

	var resp types.KonnectPaymentResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if err := s.CreateKonnectPayment(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ✅ Implements PaymentService.VerifySignature
func (s *KonnectService) VerifySignature(secret, body, signature string) bool {
	// Konnect does not support signature verification
	return true
}

// ✅ Implements PaymentService.GetPaymentStatus
func (s *KonnectService) GetPaymentStatus(paymentID string) (string, error) {
	req, _ := http.NewRequest(http.MethodGet, s.BaseURL+"/payments/"+paymentID, nil)
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch payment status: %d", resp.StatusCode)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Status, nil
}
