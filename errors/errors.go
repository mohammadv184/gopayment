package errors

type (
	ErrInternal struct {
		Message string `json:"message"`
	}
	ErrPurchaseFailed struct {
		Message string `json:"message"`
	}
	ErrInvalidPayment struct {
		Message string `json:"message"`
	}
)

func (e ErrInternal) Error() string {
	return "internal error"
}
func (e ErrPurchaseFailed) Error() string {
	return e.Message
}
func (e ErrInvalidPayment) Error() string {
	return e.Message
}
