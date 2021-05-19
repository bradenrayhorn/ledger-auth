package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
)

type HMACService struct {
	key []byte
}

func NewHMACService(key []byte) HMACService {
	return HMACService{key: key}
}

func (s HMACService) SignData(data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, s.key)
	_, err := h.Write(data)
	if err != nil {
		return []byte(""), err
	}

	return h.Sum(nil), nil
}

func (s HMACService) ValidateSignature(data []byte, signature []byte) error {
	expectedSignature, err := s.SignData(data)
	if err != nil {
		return err
	}

	if !hmac.Equal(expectedSignature, signature) {
		return errors.New("invalid signature")
	}
	return nil
}
