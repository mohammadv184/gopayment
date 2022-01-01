package test

import (
	"testing"

	"github.com/mohammadv184/gopayment/trait"
	"github.com/stretchr/testify/suite"
)

type TraitsDetailTestSuite struct {
	suite.Suite
	ExampleObj struct {
		trait.HasDetail
	}
}

func (s *TraitsDetailTestSuite) TestDetailTrait() {
	d := s.ExampleObj.GetDetail("test")
	s.Empty(d)

	s.ExampleObj.Detail("foo", "bar")
	s.ExampleObj.Detail("foo2", "bar2")
	foo := s.ExampleObj.GetDetail("foo")
	s.NotEmpty(foo)
	s.Equal("bar", foo)

	foo2 := s.ExampleObj.GetDetail("foo2")
	s.NotEmpty(foo2)
	s.Equal("bar2", foo2)

	s.Equal(false, s.ExampleObj.Has("foo3"))
	d = s.ExampleObj.GetDetail("foo3")
	s.Empty(d)

	details := s.ExampleObj.GetDetails()
	s.Equal(map[string]string{
		"foo":  "bar",
		"foo2": "bar2",
	}, details)

}
func TestTraitsDetailTestSuite(t *testing.T) {
	suite.Run(t, new(TraitsDetailTestSuite))
}
