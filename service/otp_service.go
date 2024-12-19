package service

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type OtpService interface {
	Generate() string
}

type otpService struct {
	log *zap.Logger
}

func NewOtpService(log *zap.Logger) OtpService {
	return otpService{log}
}

func (s otpService) Generate() string {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	otp := fmt.Sprintf("%06d", rng.Intn(1000000)) // Generate 6 digit OTP
	return otp
}
