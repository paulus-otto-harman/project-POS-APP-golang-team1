package validation_rules

import (
	"github.com/go-playground/validator/v10"
	"project/domain"
	"strconv"
)

func (r *Rule) PaymentStatus(fl validator.FieldLevel) bool {
	var order domain.Order
	if data, _, _, ok := fl.GetStructFieldOK2(); ok {
		order = data.Interface().(domain.Order)
	}

	if err := r.repo.Order.FindByIDOrder(&order, strconv.Itoa(int(order.ID))); err != nil {
		return false
	}

	newStatusPayment := fl.Field().Interface().(domain.StatusPayment)

	if order.StatusPayment == newStatusPayment {
		return true
	}

	switch newStatusPayment {
	case domain.OrderInProcess:
		return order.StatusPayment == ""
	case domain.OrderCancelled:
		return order.StatusPayment == domain.OrderInProcess
	case domain.OrderCompleted:
		return order.StatusPayment == domain.OrderInProcess
	default:
		return false
	}
}
