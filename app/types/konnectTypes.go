package types

import "time"

// InitKonnectPaymentRequest represents the payload for creating a payment
type InitKonnectPaymentRequest struct {
	ReceiverWalletId       string   `json:"receiverWalletId"`
	Token                  string   `json:"token"`
	Amount                 float64  `json:"amount"`
	Type                   string   `json:"type"`
	Description            string   `json:"description"`
	AcceptedPaymentMethods []string `json:"acceptedPaymentMethods"`
	Lifespan               int      `json:"lifespan"`
	CheckoutForm           bool     `json:"checkoutForm"`
	AddPaymentFeesToAmount bool     `json:"addPaymentFeesToAmount"`
	FirstName              string   `json:"firstName"`
	LastName               string   `json:"lastName"`
	PhoneNumber            string   `json:"phoneNumber"`
	Email                  string   `json:"email"`
	OrderID                string   `json:"orderId"`
	Webhook                string   `json:"webhook"`
	SilentWebhook          bool     `json:"silentWebhook,omitempty"`
	SuccessURL             string   `json:"successUrl,omitempty"`
	FailURL                string   `json:"failUrl,omitempty"`
	Theme                  string   `json:"theme"`
}

// KonnectPaymentResponse represents the API response from Konnect
type KonnectPaymentResponse struct {
	ID          int       `json:"id"`
	UserID      int       `gorm:"user_id"`
	PaymentRef  string    `json:"paymentRef" gorm:"uniqueIndex;size:100;not null"`
	PaymentLink string    `json:"payUrl"`
	Status      string    `json:"status"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// Implement generic interface
func (k *KonnectPaymentResponse) GetPaymentLink() string   { return k.PaymentLink }
func (k *KonnectPaymentResponse) GetPaymentRef() string    { return k.PaymentRef }
func (k *KonnectPaymentResponse) GetAmount() int           { return k.Amount }
func (k *KonnectPaymentResponse) GetPaymentStatus() string { return k.Status }
