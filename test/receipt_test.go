package test

import (
	"github.com/mohammadv184/gopayment/receipt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ReceiptTestSuite struct {
	suite.Suite
	Receipt *receipt.Receipt
}

func (s *ReceiptTestSuite) SetupTest() {
	s.Receipt = receipt.NewReceipt("refID", "Driver")
}
func (s *ReceiptTestSuite) TestCreateReceipt() {
	s.Equal("refID", s.Receipt.GetReferenceID())
	s.Equal("Driver", s.Receipt.GetDriver())
	s.Equal(time.Now().Format("2006-01-02"), s.Receipt.GetDate().Format("2006-01-02"))
}

func TestReceiptTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiptTestSuite))
}
