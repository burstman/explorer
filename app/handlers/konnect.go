package handlers

import (
	"encoding/json"
	"explorer/app/db"
	"explorer/app/types"
	"explorer/plugins/paymentservices"
	"io"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anthdm/superkit/kit"
)

// Konnect initializes a Konnect payment service from environment variables and starts the payment
// initialization flow.
//
// Konnect reads the following environment variables to populate a paymentservices.KonnectService:
//   - KONNECT_API_KEY:      API key for the Konnect payment provider.
//   - KONNECT_API_BASE_URL: Base URL of the Konnect API.
//   - WEBHOOK_URL:          URL to receive webhook callbacks.
//   - SUCCESS_URL:          URL to redirect to on successful payment.
//   - FAIL_URL:             URL to redirect to on failed payment.
//   - RECEIVER_WALLET_ID:   Wallet identifier that will receive funds.
//
// The constructed service is passed, together with the provided kit, to KonnectInitPayment.
// The kit parameter must be a non-nil *kit.Kit and is used by KonnectInitPayment for any
// application-specific context or resources required during initialization.
//
// Returns any error returned by KonnectInitPayment. This function's side effects are limited
// to reading environment variables and invoking KonnectInitPayment; it does not itself
// perform network calls beyond what KonnectInitPayment may perform.
func Konnect(kit *kit.Kit) error {
	service := &paymentservices.KonnectService{
		APIKey:     os.Getenv("KONNECT_API_KEY"),
		BaseURL:    os.Getenv("KONNECT_API_BASE_URL"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),
		SuccessURL: os.Getenv("SUCCESS_URL"),
		FailURL:    os.Getenv("FAIL_URL"),
		WalletID:   os.Getenv("RECEIVER_WALLET_ID"),
	}
	return KonnectInitPayment(kit, service)

}

func KonnectInitPayment(kit *kit.Kit, service types.PaymentService) error {
	authuser := kit.Auth().(types.AuthUser)
	userID := authuser.GetUserID()

	err := kit.Request.ParseForm()
	if err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	//strAmount := kit.Request.PostFormValue("amount")

	//amount, err := strconv.ParseFloat(strAmount, 64)

	// if err != nil {
	// 	return fmt.Errorf("failed to parse form: %v", err)
	// }

	var user types.User
	err = db.Get().First(&user, userID).Error
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Build request

	konnectservice, ok := service.(*paymentservices.KonnectService)
	if !ok {
		return fmt.Errorf("invalid konnect payment service type")
	}
	log.Println("webhook url:", konnectservice.WebhookURL)

	req := &types.InitKonnectPaymentRequest{
		ReceiverWalletId: konnectservice.WalletID,
		Type:             "immediate",
		Token:            "TND",
		//Amount:           amount * 1000,
		Amount:      1000,
		Description: "Test Payment for order #1234",
		AcceptedPaymentMethods: []string{
			"wallet",
			"bank_card",
			"e-DINAR",
			"konnect",
		},
		CheckoutForm:           false,
		AddPaymentFeesToAmount: true,
		Lifespan:               10,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		Email:                  user.Email,
		PhoneNumber:            user.PhoneNumber,
		OrderID:                "123456",
		Webhook:                konnectservice.WebhookURL,
		//SilentWebhook:          false,
		Theme: "light",
	}

	// Call payment service
	resp, err := service.InitPayment(int(userID), req)
	if err != nil {
		return fmt.Errorf("payment initialization failed: %v", err)
	}

	log.Println("Redirecting to payment link:", resp.GetPaymentLink())

	return kit.Redirect(http.StatusSeeOther, resp.GetPaymentLink())
}

func HandlePaymentStatus(kit *kit.Kit) error {
	log.Println("Handling payment status request")
	authUser := kit.Auth().(types.AuthUser)
	userID := authUser.GetUserID()

	// Parse form data
	if err := kit.Request.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Fetch user payments
	payments := GetPaymentRef(int(userID))
	if len(payments) == 0 {
		return fmt.Errorf("no payments found for user %d", userID)
	}

	// Example: use the most recent payment (or adjust logic as needed)
	payment := payments[len(payments)-1]

	baseURL := os.Getenv("KONNECT_API_BASE_URL")
	if baseURL == "" {
		return fmt.Errorf("KONNECT_API_BASE_URL not set in environment")
	}

	// Construct full endpoint: GET /payments/:paymentId
	url := fmt.Sprintf("%s/payments/%s", baseURL, payment)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Optional: set headers (e.g., auth token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call Konnect API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Konnect API error: %s", string(body))
	}

	var paymentStatus map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&paymentStatus); err != nil {
		return fmt.Errorf("failed to decode Konnect API response: %v", err)
	}

	log.Printf("Payment status for user %d: %+v", userID, paymentStatus)

	return kit.JSON(http.StatusOK, paymentStatus)
}

func GetPaymentRef(userID int) []string {
	var payments []types.KonnectPaymentResponse
	if err := db.Get().Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		log.Printf("failed to retrieve payments: %v", err)
		return nil
	}
	var paymentRefs []string
	for _, p := range payments {
		paymentRefs = append(paymentRefs, p.PaymentRef)
	}
	return paymentRefs
}
