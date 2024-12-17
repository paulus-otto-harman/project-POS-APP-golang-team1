package domain

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type StatusPayment string

const (
	OrderInProcess StatusPayment = "In Process"
	OrderCompleted StatusPayment = "Completed"
	OrderCancelled StatusPayment = "Cancelled"
)

type StatusKitchen string

const (
	OrderInTheKitchen StatusKitchen = "In The Kitchen"
	OrderCookingNow   StatusKitchen = "Cooking Now"
	OrderReadyToServe StatusKitchen = "Ready To Serve"
)

type Order struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id" swaggerignore:"true"`
	TableID         uint           `gorm:"not null" json:"table_id"`
	Table           Table          `gorm:"foreignKey:TableID;references:ID"`
	Name            string         `gorm:"size:100" json:"name"`
	CodeOrder       string         `gorm:"size:50;unique" json:"code_order"`
	Tax             float64        `gorm:"type:decimal(4,2);not null;default:10.0" json:"tax"`
	PaymentMethodID uint           `gorm:"default:null" json:"payment_method_id"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodID;references:ID"`
	StatusPayment   StatusPayment  `gorm:"type:status_payment;default:'In Process'" json:"status_payment" example:"In Process"`
	StatusKitchen   StatusKitchen  `gorm:"type:status_kitchen;default:'In The Kitchen'" json:"status_kitchen" example:"In The Kitchen"`
	OrderItems      []OrderItem    `gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

func generateCodeOrder(db *gorm.DB) (string, error) {
	var lastOrder Order
	err := db.Order("id desc").First(&lastOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", fmt.Errorf("failed to retrieve last order: %v", err)
	}

	lastCode := "ORD0000"
	if lastOrder.CodeOrder != "" {
		lastCode = lastOrder.CodeOrder
	}

	codeNum, err := strconv.Atoi(lastCode[3:])
	if err != nil {
		return "", fmt.Errorf("failed to parse last order code number: %v", err)
	}

	newCodeNum := codeNum + 1
	newCode := fmt.Sprintf("ORD%04d", newCodeNum)

	return newCode, nil
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	codeOrder, err := generateCodeOrder(tx)
	if err != nil {
		return fmt.Errorf("failed to generate code_order: %v", err)
	}
	o.CodeOrder = codeOrder
	return nil
}

// func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
// 	var totalSubTotal float64

// 	for i, item := range o.OrderItems {
// 		var product Product
// 		if err := tx.First(&product, item.ProductID).Error; err != nil {
// 			return fmt.Errorf("product not found for product_id %d", item.ProductID)
// 		}
// 		o.OrderItems[i].SubTotal = product.Price * float64(item.Quantity)
// 		totalSubTotal += o.OrderItems[i].SubTotal
// 	}

// 	o.Amount = totalSubTotal + (totalSubTotal * o.Tax / 100)
// 	return nil
// }

type OrderItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id" swaggerignore:"true"`
	OrderID   uint      `gorm:"not null" json:"order_id" example:"1"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID"`
	ProductID uint      `gorm:"not null" json:"product_id" example:"1"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID"`
	Quantity  int       `gorm:"not null" json:"quantity" example:"2"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

// func (oi *OrderItem) BeforeSave(tx *gorm.DB) (err error) {

// 	if oi.Quantity <= 0 {
// 		return errors.New("quantity must be greater than 0")
// 	}
// 	var product Product
// 	if err := tx.First(&product, oi.ProductID).Error; err != nil {
// 		return errors.New("product not found")
// 	}

// 	oi.SubTotal = product.Price * float64(oi.Quantity)
// 	return nil
// }

func OrderSeed() []Order {
	return []Order{
		{
			ID:      1,
			TableID: 1,
			Name:    "John Doe",
			// CodeOrder:       "ORD001",
			Tax:             10.0,
			PaymentMethodID: 1,
			StatusPayment:   "In Process",
			OrderItems: []OrderItem{
				{
					// OrderID: 1,
					ProductID: 1,
					Quantity:  2,
				},
				{
					// OrderID: 1,
					ProductID: 2,
					Quantity:  1,
				},
			},
		},
		// {
		// 	ID: 2,
		// 	TableID:         2,
		// 	Name:            "Alex",
		// 	// CodeOrder:       "ORD002",
		// 	Tax:             10.0,
		// 	PaymentMethodID: 2,
		// 	StatusPayment:   "Completed",
		// 	OrderItems: []OrderItem{
		// 		{
		// 			OrderID: 2,
		// 			ProductID: 3,
		// 			Quantity:  3,
		// 		},
		// 		{
		// 			OrderID: 2,
		// 			ProductID: 12,
		// 			Quantity:  3,
		// 		},
		// 		{
		// 			OrderID: 2,
		// 			ProductID: 6,
		// 			Quantity:  3,
		// 		},
		// 	},
		// },
		// {
		// 	ID: 3,
		// 	TableID:         3,
		// 	Name:            "Elia",
		// 	// CodeOrder:       "ORD003",
		// 	Tax:             10.0,
		// 	PaymentMethodID: 3,
		// 	StatusPayment:   "Cancelled",
		// 	OrderItems: []OrderItem{
		// 		{
		// 			OrderID: 3,
		// 			ProductID: 4,
		// 			Quantity:  1,
		// 		},
		// 		{
		// 			OrderID: 3,
		// 			ProductID: 15,
		// 			Quantity:  1,
		// 		},
		// 	},
		// },
		// {
		// 	ID: 4,
		// 	TableID:         4,
		// 	Name:            "Smith",
		// 	// CodeOrder:       "ORD004",
		// 	Tax:             10.0,
		// 	PaymentMethodID: 2,
		// 	StatusPayment:   "In Process",
		// 	OrderItems: []OrderItem{
		// 		{
		// 			OrderID: 4,
		// 			ProductID: 5,
		// 			Quantity:  5,
		// 		},
		// 		{
		// 			OrderID: 4,
		// 			ProductID: 17,
		// 			Quantity:  5,
		// 		},
		// 		{
		// 			OrderID: 4,
		// 			ProductID: 2,
		// 			Quantity:  5,
		// 		},
		// 	},
		// },
		// {
		// 	ID: 5,
		// 	TableID:         5,
		// 	Name:            "Bob Brown",
		// 	// CodeOrder:       "ORD005",
		// 	Tax:             10.0,
		// 	PaymentMethodID: 0,
		// 	StatusPayment:   "Completed",
		// 	OrderItems: []OrderItem{
		// 		{
		// 			OrderID: 5,
		// 			ProductID: 5,
		// 			Quantity:  5,
		// 		},
		// 	},
		// },
	}

}
