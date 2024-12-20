package validation_rules

import (
	"go.uber.org/zap"
	"project/repository"
)

type Rule struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewValidation(repo repository.Repository, log *zap.Logger) struct{ Rule Rule } {
	return struct{ Rule Rule }{Rule: Rule{repo, log}}
}
