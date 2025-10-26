package handlers

import (
	"explorer/app/db"
	"explorer/plugins/paymentservices"
	"log"

	"explorer/app/types"
	"fmt"
	"net/http"

	"github.com/anthdm/superkit/kit"
)

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
		},
		CheckoutForm:           false,
		AddPaymentFeesToAmount: true,
		Lifespan:               15,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		Email:                  user.Email,
		PhoneNumber:            user.PhoneNumber,
		OrderID:                "123456",
		Webhook:                konnectservice.WebhookURL,
		SilentWebhook:          false,
		SuccessURL:             konnectservice.SuccessURL,
		FailURL:                konnectservice.FailURL,
		Theme:                  "light",
	}

	// Call payment service
	resp, err := service.InitPayment(req)
	if err != nil {
		return fmt.Errorf("payment initialization failed: %v", err)
	}

	log.Println("Redirecting to payment link:", resp.GetPaymentLink())

	return kit.Redirect(http.StatusSeeOther, resp.GetPaymentLink())
}
