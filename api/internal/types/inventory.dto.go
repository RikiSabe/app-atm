package types

type InventoryItemDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	NetPrice  float64 `json:"net_price"`
	SalePrice float64 `json:"sale_price"`
}

type InventoryItemInputDTO struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	NetPrice  float64 `json:"net_price"`
	SalePrice float64 `json:"sale_price"`
}

type SellableInventoryItemDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	SalePrice float64 `json:"sale_price"`
}

type ConsumeInventoryDTO struct {
	CheckoutID string                    `json:"checkout_id"`
	Items      []InventoryConsumptionDTO `json:"items"`
}

type InventoryConsumptionDTO struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
