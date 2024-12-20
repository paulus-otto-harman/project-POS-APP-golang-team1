package service

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"project/helper/validation_rules"
	"project/repository"
)

type ValidationService interface {
}

type validationService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewValidationService(repo repository.Repository, log *zap.Logger) ValidationService {
	validation := validation_rules.NewValidation(repo, log)

	if customValidator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidator.RegisterValidation("paymentStatus", validation.Rule.PaymentStatus)
	}

	return &validationService{repo, log}
}
