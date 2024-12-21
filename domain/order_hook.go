package domain

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"
)

// Utility Functions
func generateCodeOrder(db *gorm.DB) (string, error) {
	var lastOrder Order
	err := db.Unscoped().Order("id desc").First(&lastOrder).Error
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

func validateTable(tx *gorm.DB, tableID uint) error {
	var table Table

	if err := tx.First(&table, tableID).Error; err != nil {
		return fmt.Errorf("failed to retrieve table: %v", err)
	}

	if !table.Status {
		return fmt.Errorf(table.Name + " is already reserved")
	}

	return nil
}
func changeTable(tx *gorm.DB, orderID, tableID uint) error {
	var oldOrder Order
	if err := tx.Unscoped().First(&oldOrder, orderID).Error; err != nil {
		return fmt.Errorf("failed to retrieve old order: %v", err)
	}

	if oldOrder.TableID != tableID {
		log.Println("masuk before update oldtable != table")

		if err := validateTable(tx, tableID); err != nil {
			return err
		}
		if err := updateTableStatus(tx, oldOrder.TableID, true); err != nil {
			return err
		}

		if err := updateTableStatus(tx, tableID, false); err != nil {
			return err
		}

	}

	return nil
}
func restoreStockForCancelledOrder(tx *gorm.DB, orderID uint) error {
	var oldOrder Order
	if err := tx.Unscoped().First(&oldOrder, orderID).Error; err != nil {
		return fmt.Errorf("failed to retrieve old order: %v", err)
	}

	var orderItems []OrderItem
	if err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return fmt.Errorf("failed to retrieve order items for cancelled order: %v", err)
	}

	if oldOrder.StatusPayment != OrderCancelled {
		for _, item := range orderItems {
			if err := adjustProductStock(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
	}

	return nil
}
func overWriteOrderItem(tx *gorm.DB, orderID uint, orderItems []OrderItem) error {
	var existingItems []OrderItem
	if err := tx.Where("order_id = ?", orderID).Find(&existingItems).Error; err != nil {
		return fmt.Errorf("failed to retrieve existing order items: %v", err)
	}

	newItemIDs := make(map[uint]bool)
	for _, item := range orderItems {
		newItemIDs[item.ID] = true
	}

	for _, existingItem := range existingItems {
		if !newItemIDs[existingItem.ID] {
			if err := tx.Delete(&existingItem).Error; err != nil {
				return fmt.Errorf("failed to delete removed order item: %v", err)
			}
		}
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

// Hook Order
func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	if o.ID == 0 {
		codeOrder, err := generateCodeOrder(tx)
		if err != nil {
			return fmt.Errorf("failed to generate code_order: %v", err)
		}
		o.CodeOrder = codeOrder
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {

	if err := validateTable(tx, o.TableID); err != nil {
		return err
	}
	if err := updateTableStatus(tx, o.TableID, false); err != nil {
		return err
	}
	return nil
}

func (o *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	if o.PaymentMethodID != nil {
		o.StatusPayment = OrderCompleted
		o.StatusKitchen = OrderReadyToServe
	}

	var order Order
	if err := tx.First(&order, o.ID).Error; err != nil {
		return fmt.Errorf("failed to retrieve table: %v", err)
	}

	if order.StatusPayment == OrderInProcess && o.PaymentMethodID == nil {
		log.Println(order.StatusPayment, "masuk before update")

		if err := changeTable(tx, o.ID, o.TableID); err != nil {
			return err
		}
	}

	if o.StatusPayment == OrderCancelled {
		if err := restoreStockForCancelledOrder(tx, o.ID); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) AfterUpdate(tx *gorm.DB) (err error) {
	var order Order
	if err := tx.First(&order, o.ID).Error; err != nil {
		return fmt.Errorf("failed to retrieve table: %v", err)
	}
	if order.StatusPayment != OrderInProcess {
		log.Println(order.StatusPayment, "masuk after update")

		if err := updateTableStatus(tx, o.TableID, true); err != nil {
			return fmt.Errorf("failed to update table status: %v", err)
		}
	}

	if err := overWriteOrderItem(tx, o.ID, o.OrderItems); err != nil {
		return err
	}
  
	return nil
}

func (o *Order) AfterDelete(tx *gorm.DB) (err error) {
	var orderItems []OrderItem
	if err := tx.Where("order_id = ?", o.ID).Find(&orderItems).Error; err != nil {
		return fmt.Errorf("failed to retrieve order items for deleted order: %v", err)
	}

	for _, item := range orderItems {
		if err := adjustProductStock(tx, item.ProductID, item.Quantity); err != nil {
			return fmt.Errorf("failed to restore stock for product ID %d: %v", item.ProductID, err)
		}
	}

	if o.StatusPayment == OrderInProcess {
		log.Println(o.StatusPayment, "masuk after delete")
		if err := updateTableStatus(tx, o.TableID, true); err != nil {
			return fmt.Errorf("failed to update table status: %v", err)
		}
	}
	return nil
}

// Hook Order Item
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
