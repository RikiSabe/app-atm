package models

import "time"

type InventoryItem struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name;size:120;not null" json:"name"`
	Quantity  int       `gorm:"column:quantity;not null" json:"quantity"`
	NetPrice  float64   `gorm:"column:net_price;type:numeric(12,2);not null;default:0" json:"net_price"`
	SalePrice float64   `gorm:"column:sale_price;type:numeric(12,2);not null;default:0" json:"sale_price"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
}

func (InventoryItem) TableName() string {
	return "inventory_items"
}

type InventoryClosure struct {
	ID        string    `gorm:"column:id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (InventoryClosure) TableName() string {
	return "inventory_closures"
}
