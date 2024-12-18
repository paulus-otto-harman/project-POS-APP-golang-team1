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
	ID              uint           `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	TableID         uint           `gorm:"not null" json:"table_id" example:"1"`
	Table           Table          `gorm:"foreignKey:TableID;references:ID"`
	Name            string         `gorm:"size:100" json:"name"`
	CodeOrder       string         `gorm:"size:50;unique" json:"code_order"`
	Tax             float64        `gorm:"type:decimal(4,2);not null;default:10.0" json:"tax"`
	PaymentMethodID *uint          `gorm:"default:null" json:"payment_method_id" example:"1"`
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

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	var table Table

	if err := tx.First(&table, o.TableID).Error; err != nil {
		return fmt.Errorf("failed to retrieve table: %v", err)
	}

	if !table.Status {
		return fmt.Errorf(table.Name + " is already reserved")
	}

	return nil
}

func (o *Order) AfterCreate(tx *gorm.DB) (err error) {

	if err := updateTableStatus(tx, o.TableID, false); err != nil {
		return fmt.Errorf("failed to update table status: %v", err)
	}
	return nil
}

func (o *Order) AfterUpdate(tx *gorm.DB) (err error) {
	if o.StatusPayment == OrderCompleted {

		if err := updateTableStatus(tx, o.TableID, true); err != nil {
			return fmt.Errorf("failed to update table status: %v", err)
		}
	}
	return nil
}

func updateTableStatus(tx *gorm.DB, tableID uint, status bool) error {
	var table Table
	if err := tx.First(&table, tableID).Error; err != nil {
		return fmt.Errorf("failed to retrieve table: %v", err)
	}

	table.Status = status
	if err := tx.Save(&table).Error; err != nil {
		return fmt.Errorf("failed to update table status: %v", err)
	}
	return nil
}

type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	OrderID   uint      `gorm:"not null" json:"order_id" example:"1" swaggerignore:"true"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID" swaggerignore:"true"`
	ProductID uint      `gorm:"not null" json:"product_id" example:"1"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" swaggerignore:"true"`
	Quantity  int       `gorm:"not null" json:"quantity" example:"2"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

func (oi *OrderItem) AfterCreate(tx *gorm.DB) (err error) {
	if err := adjustProductStock(tx, oi.ProductID, -oi.Quantity); err != nil {
		return fmt.Errorf("failed: %v", err)
	}
	return nil
}

func (oi *OrderItem) AfterUpdate(tx *gorm.DB) (err error) {
	var old OrderItem
	if err := tx.First(&old, oi.ID).Error; err != nil {
		return fmt.Errorf("failed to retrieve old order item: %v", err)
	}

	stockDifference := oi.Quantity - old.Quantity
	if err := adjustProductStock(tx, oi.ProductID, -stockDifference); err != nil {
		return fmt.Errorf("failed to adjust product stock: %v", err)
	}
	return nil
}

func (oi *OrderItem) AfterDelete(tx *gorm.DB) (err error) {
	if err := adjustProductStock(tx, oi.ProductID, oi.Quantity); err != nil {
		return fmt.Errorf("failed to restore product stock: %v", err)
	}
	return nil
}

func adjustProductStock(tx *gorm.DB, productID uint, quantityChange int) error {
	var product Product
	if err := tx.First(&product, productID).Error; err != nil {
		return fmt.Errorf("failed to retrieve product: %v", err)
	}

	newStock := product.Stock + quantityChange
	if newStock < 0 {
		return fmt.Errorf("insufficient stock for product %s", product.Name)
	}

	product.Stock = newStock
	if err := tx.Save(&product).Error; err != nil {
		return fmt.Errorf("failed to update product stock: %v", err)
	}
	return nil
}

func OrderSeed() []Order {
	return []Order{
		{
			TableID: 1,
			Name:    "John Doe",
			OrderItems: []OrderItem{
				{
					ProductID: 2,
					Quantity:  2,
				},
			},
		},
	}
}
