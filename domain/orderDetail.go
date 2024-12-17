package domain

import "encoding/json"

type OrderDetail struct {
	OrderID           int             `json:"order_id"`
	Name              string          `json:"name"`
	TableName         string          `json:"table_name"`
	PaymentMethodName string          `json:"payment_method_name"`
	OrderItems        json.RawMessage `json:"order_items"`
	Total             float64         `json:"total"`
}
