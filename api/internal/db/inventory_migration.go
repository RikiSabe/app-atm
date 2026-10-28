package db

import (
	"fmt"

	"app-atm-api/internal/models"

	"gorm.io/gorm"
)

func MigrateInventory() error {
	return GDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.InventoryItem{}, &models.InventoryClosure{}); err != nil {
			return fmt.Errorf("actualizar el esquema de inventario: %w", err)
		}

		if tx.Migrator().HasColumn(&models.InventoryItem{}, "price") {
			if err := tx.Exec(
				"UPDATE inventory_items SET net_price = price, sale_price = price",
			).Error; err != nil {
				return fmt.Errorf("copiar precios anteriores del inventario: %w", err)
			}
			if err := tx.Migrator().DropColumn(&models.InventoryItem{}, "price"); err != nil {
				return fmt.Errorf("eliminar la columna anterior de precio: %w", err)
			}
		}
		return nil
	})
}
