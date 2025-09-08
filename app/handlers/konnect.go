package handlers

import (
	"encoding/json"
	"explorer/app/db"
	"explorer/plugins/paymentservices"

	"explorer/app/types"
	"fmt"
	"net/http"
	"os"

	"github.com/anthdm/superkit/kit"
)

func KonnectInitPayment(kit *kit.Kit) error {
	authuser := kit.Auth().(types.AuthUser)
	userID := authuser.GetUserID()

	var user types.User
	err := db.Get().First(&user, userID).Error
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Create service instance
	service := &paymentservices.KonnectService{
		APIKey:     os.Getenv("KONNECT_API_KEY"),
		BaseURL:    os.Getenv("KONNECT_API_BASE_URL"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),
		SuccessURL: os.Getenv("SUCCESS_URL"),
		FailURL:    os.Getenv("FAIL_URL"),
		WalletID:   os.Getenv("RECEIVER_WALLET_ID"),
	}

	// Build request
	req := &types.InitKonnectPaymentRequest{
		ReceiverWalletId: service.WalletID,
		Type:             "immediate",
		Amount:           1000,
		Description:      "Test Payment for order #1234",
		AcceptedPaymentMethods: []string{
			"card",
			"bank_transfer",
		},
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Webhook:       service.WebhookURL,
		SilentWebhook: true,
		SuccessURL:    service.SuccessURL,
		FailURL:       service.FailURL,
		Theme:         "light",
	}

	// Call payment service
	resp, err := service.InitPayment(req)
	if err != nil {
		http.Error(kit.Response, fmt.Sprintf("payment initialization failed: %v", err), http.StatusBadRequest)
		return nil
	}

	kit.Response.Header().Set("Content-Type", "application/json")
	kit.Response.WriteHeader(http.StatusCreated)
	return json.NewEncoder(kit.Response).Encode(resp)
}
