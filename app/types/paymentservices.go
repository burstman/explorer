package types

// Generic interface input for any payment service
type InitPaymentRequest interface{}

// Generic interface output from any payment service
type PaymentResponse interface {
	GetPaymentRef() string
	GetAmount() int
	GetStatus() string
	GetPaymentLink() string
}

type PaymentService interface {
	InitPayment(InitPaymentRequest) (PaymentResponse, error)
	VerifySignature(secret, body, signature string) bool // optional if provider supports it
	GetPaymentStatus(paymentID string) (string, error)   // optional
}

func NewOnlinePaymentService(paymentService PaymentService) PaymentService {
	return paymentService
}
