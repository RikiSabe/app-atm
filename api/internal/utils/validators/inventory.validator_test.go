package validators

import (
	"strings"
	"testing"

	"app-atm-api/internal/types"
)

func TestValidateItem(t *testing.T) {
	tests := []struct {
		name      string
		input     types.InventoryItemInputDTO
		wantValid bool
	}{
		{name: "valid item trims name", input: types.InventoryItemInputDTO{Name: "  Bebida  ", Quantity: 0, NetPrice: 0, SalePrice: 0}, wantValid: true},
		{name: "maximum database price", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: 1, NetPrice: maxInventoryPrice, SalePrice: maxInventoryPrice}, wantValid: true},
		{name: "empty name", input: types.InventoryItemInputDTO{Quantity: 1, NetPrice: 1, SalePrice: 1}},
		{name: "negative quantity", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: -1, NetPrice: 1, SalePrice: 1}},
		{name: "negative net price", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: 1, NetPrice: -0.01, SalePrice: 1}},
		{name: "negative sale price", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: 1, NetPrice: 1, SalePrice: -0.01}},
		{name: "net price with fractional cents", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: 1, NetPrice: 1.230001, SalePrice: 1}},
		{name: "sale price with excess precision", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: 1, NetPrice: 1, SalePrice: 1.234}},
		{name: "net price exceeds database precision", input: types.InventoryItemInputDTO{Name: "Bebida", Quantity: 1, NetPrice: 10000000000, SalePrice: 1}},
		{name: "long name", input: types.InventoryItemInputDTO{Name: strings.Repeat("a", 121), Quantity: 1, NetPrice: 1, SalePrice: 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := NewInventoryValidator()
			gotValid := validator.ValidateItem(&test.input)
			if gotValid != test.wantValid {
				t.Errorf("ValidateItem() = %v, want %v; errors: %v", gotValid, test.wantValid, validator.GetErrors())
			}
			if test.wantValid && test.input.Name != "Bebida" {
				t.Errorf("ValidateItem() did not trim name: got %q", test.input.Name)
			}
		})
	}
}
