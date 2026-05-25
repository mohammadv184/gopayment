package gopayment

import "errors"

var (
	ErrGatewayAlreadyRegistered = errors.New("gopayment: gateway already registered")
	ErrGatewayNotFound          = errors.New("gopayment: gateway not found")
	ErrGatewayNotRedirector     = errors.New("gopayment: gateway is not a redirector")
	ErrGatewayNotRefunder       = errors.New("gopayment: gateway is not a refunder")
)
