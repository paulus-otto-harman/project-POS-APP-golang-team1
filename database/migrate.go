package database

import (
	"project/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	var err error

	if err = dropTables(db); err != nil {
		return err
	}

	if err = setupJoinTables(db); err != nil {
		return err
	}

	if err = autoMigrates(db); err != nil {
		return err
	}

	return createViews(db)
}

func autoMigrates(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Permission{},
		&domain.User{},
		&domain.Reservation{},
		&domain.Notification{},
		&domain.Category{},
		&domain.Product{},
		&domain.Table{},
		&domain.PaymentMethod{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Profile{},
		&domain.PasswordResetToken{},
	)
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&domain.PasswordResetToken{},
		&domain.Profile{},
		&domain.User{},
		&domain.Reservation{},
		&domain.Notification{},
		&domain.Category{},
		&domain.Product{},
		&domain.Table{},
		&domain.PaymentMethod{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Permission{},
		"user_permissions",
		&domain.UserNotification{},
	)
}

func setupJoinTables(db *gorm.DB) error {
	var err error
	if err = db.SetupJoinTable(&domain.User{}, "Permissions", &domain.UserPermission{}); err != nil {
		return err
	}

	if err = db.SetupJoinTable(&domain.User{}, "Notifications", &domain.UserNotification{}); err != nil {
		return err
	}
	return err
}

func createViews(db *gorm.DB) error {
	var err error
	if err = queryOrderDetail(db); err != nil {
		return err
	}
	return err
}

func queryOrderDetail(db *gorm.DB) error {

	query := db.Raw(`
	SELECT 
    o.id AS order_id, o.name, o.code_order, o.status_payment, o.status_kitchen, t.name AS table_name, pm.name AS payment_method_name,
	to_char(o.created_at, 'FMDay, DD-Mon-YYYY') as date_order,
    to_char(o.created_at, 'FMHH12:MI AM') as time_order,
   	jsonb_agg(
        jsonb_build_object(
            'order_item_id', oi.id,
            'product_name', p.name,
            'product_price', p.price,
            'quantity', oi.quantity,
            'sub_total', (oi.quantity * p.price)
        )
    ) AS order_items,
    SUM(oi.quantity * p.price) as total
	FROM orders o
	LEFT JOIN tables t ON o.table_id = t.id
	LEFT JOIN payment_methods pm ON o.payment_method_id = pm.id
	LEFT JOIN order_items oi ON o.id = oi.order_id
	LEFT JOIN products p ON oi.product_id = p.id
	WHERE o.deleted_at IS NULL
	GROUP BY 
    o.id, t.id, t.name, pm.id, pm.name
   	ORDER BY o.id;
	`)

	return db.Migrator().CreateView("order_details", gorm.ViewOption{Query: query, Replace: true})
}
