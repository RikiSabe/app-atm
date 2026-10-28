package validators

import (
	"math"
	"strings"

	"app-atm-api/internal/types"
)

type InventoryValidator struct {
	Errors map[string]string
}

const maxInventoryPrice = 9999999999.99

func NewInventoryValidator() *InventoryValidator {
	return &InventoryValidator{Errors: make(map[string]string)}
}

func (validator *InventoryValidator) ValidateItem(item *types.InventoryItemInputDTO) bool {
	validator.Errors = make(map[string]string)

	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		validator.Errors["name"] = "El nombre del artículo es obligatorio."
	} else if len([]rune(item.Name)) > 120 {
		validator.Errors["name"] = "El nombre no puede exceder 120 caracteres."
	}

	if item.Quantity < 0 {
		validator.Errors["quantity"] = "La cantidad debe ser un entero igual o mayor que cero."
	}

	validatePrice := func(field string, price float64) {
		precisionTolerance := math.Max(0.00000001, math.Abs(price)*0.00000000000002)
		if math.IsNaN(price) || math.IsInf(price, 0) || price < 0 ||
			price > maxInventoryPrice || math.Abs(price*100-math.Round(price*100)) > precisionTolerance {
			validator.Errors[field] = "El precio debe ser válido, igual o mayor que cero y tener como máximo dos decimales."
		}
	}
	validatePrice("net_price", item.NetPrice)
	validatePrice("sale_price", item.SalePrice)

	return len(validator.Errors) == 0
}

func (validator *InventoryValidator) GetErrors() map[string]string {
	return validator.Errors
}
