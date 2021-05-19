package tests

import (
	"testing"

	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/stretchr/testify/suite"
)

type HMACSuite struct {
	suite.Suite
	svc services.HMACService
}

func (s *HMACSuite) SetupTest() {
	s.svc = services.NewHMACService([]byte("super secret key"))
}

func (s *HMACSuite) TestCanSignDataAndVerifyData() {
	res, err := s.svc.SignData([]byte("my message"))
	s.Require().Nil(err)

	err = s.svc.ValidateSignature([]byte("my message"), res)
	s.Require().Nil(err)
}

func (s *HMACSuite) TestCanSignDataAndWillNotAcceptWrongMessage() {
	res, err := s.svc.SignData([]byte("my message"))
	s.Require().Nil(err)

	err = s.svc.ValidateSignature([]byte("my message2"), res)
	s.Require().NotNil(err)
}

func (s *HMACSuite) TestCanSignDataAndWillNotAcceptWrongSignature() {
	err := s.svc.ValidateSignature([]byte("my message"), []byte("wrong signature"))
	s.Require().NotNil(err)
}

func TestHMACSuite(t *testing.T) {
	suite.Run(t, new(HMACSuite))
}
