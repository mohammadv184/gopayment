package test

import (
	"github.com/mohammadv184/gopayment/traits"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TraitsDetailTestSuite struct {
	suite.Suite
	ExampleObj struct {
		traits.HasDetail
	}
}

func (s *TraitsDetailTestSuite) TestDetailTrait() {
	_, err := s.ExampleObj.GetDetail("test")
	s.NotNil(err)

	s.ExampleObj.Detail("foo", "bar")
	s.ExampleObj.Detail("foo2", "bar2")
	foo, err := s.ExampleObj.GetDetail("foo")
	s.Nil(err)
	s.Equal("bar", foo)

	foo2, err := s.ExampleObj.GetDetail("foo2")
	s.Nil(err)
	s.Equal("bar2", foo2)

	_, err = s.ExampleObj.GetDetail("foo3")
	s.NotNil(err)

	details := s.ExampleObj.GetDetails()
	s.Equal(map[string]string{
		"foo":  "bar",
		"foo2": "bar2",
	}, details)

}
func TestTraitsDetailTestSuite(t *testing.T) {
	suite.Run(t, new(TraitsDetailTestSuite))
}
