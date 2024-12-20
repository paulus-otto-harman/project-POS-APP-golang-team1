package seeder

import "project/domain"

func PaymentMethodSeed() []domain.PaymentMethod {
	return []domain.PaymentMethod{
		{
			Name: "Cash",
		},
		{
			Name: "Credit Card",
		},
		{
			Name: "E-Wallet",
		},
	}
}