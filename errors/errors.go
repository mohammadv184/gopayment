package errors

type (
	// ErrInternal represents an internal error.
	ErrInternal struct {
		Message string `json:"message"`
	}
	// ErrPurchaseFailed represents a purchase failed error.
	ErrPurchaseFailed struct {
		Message string `json:"message"`
	}
	// ErrInvalidPayment represents an invalid payment error.
	ErrInvalidPayment struct {
		Message string `json:"message"`
	}
)

// Error returns the error message.
func (e ErrInternal) Error() string {
	return e.Message
}

// Error returns the error message.
func (e ErrPurchaseFailed) Error() string {
	return e.Message
}

// Error returns the error message.
func (e ErrInvalidPayment) Error() string {
	return e.Message
}
