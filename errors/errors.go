package errors

type (
	ErrInternal       struct{}
	ErrPurchaseFailed struct{}
	ErrInvalidPayment struct{}
)

func (e ErrInternal) Error() string {
	return "internal error"
}
func (e ErrPurchaseFailed) Error() string {
	return "purchase failed"
}
func (e ErrInvalidPayment) Error() string {
	return "invalid payment"
}
