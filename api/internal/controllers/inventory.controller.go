package controllers

import (
	"errors"
	"log"
	"net/http"

	"app-atm-api/internal/db"
	"app-atm-api/internal/models"
	"app-atm-api/internal/types"
	tools "app-atm-api/internal/utils"
	"app-atm-api/internal/utils/validators"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var input types.InventoryItemInputDTO
	if !tools.DecodeJSON(w, r, &input) {
		return
	}
	validator := validators.NewInventoryValidator()
	if !tools.ValidateAndRespond(w, validator.ValidateItem(&input), validator) {
		return
	}

	item := models.InventoryItem{
		ID:        uuid.NewString(),
		Name:      input.Name,
		Quantity:  input.Quantity,
		NetPrice:  input.NetPrice,
		SalePrice: input.SalePrice,
	}
	if err := db.GDB.Create(&item).Error; err != nil {
		log.Printf("No se pudo crear el artículo de inventario: %v", err)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo guardar el artículo."})
		return
	}

	tools.JsonResponse(w, http.StatusCreated, toInventoryItemDTO(item))
}

func GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	var items []models.InventoryItem
	if err := db.GDB.Order("name ASC").Find(&items).Error; err != nil {
		log.Printf("No se pudo consultar el inventario: %v", err)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo cargar el inventario."})
		return
	}

	result := make([]types.InventoryItemDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toInventoryItemDTO(item))
	}
	tools.EncodeJSON(w, result)
}

func GetSellableInventoryItems(w http.ResponseWriter, r *http.Request) {
	var items []types.SellableInventoryItemDTO
	if err := db.GDB.Model(&models.InventoryItem{}).
		Select("id, name, quantity, sale_price").
		Where("quantity > 0").
		Order("name ASC").
		Scan(&items).Error; err != nil {
		log.Printf("No se pudo consultar los artículos para venta: %v", err)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo cargar el inventario."})
		return
	}
	if items == nil {
		items = []types.SellableInventoryItemDTO{}
	}
	tools.EncodeJSON(w, items)
}

func ConsumeInventoryItems(w http.ResponseWriter, r *http.Request) {
	var input types.ConsumeInventoryDTO
	if !tools.DecodeJSON(w, r, &input) {
		return
	}
	if _, err := uuid.Parse(input.CheckoutID); err != nil {
		tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El ID del cierre no es válido."})
		return
	}
	if len(input.Items) == 0 {
		tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El cierre no contiene productos para descontar."})
		return
	}

	seen := make(map[string]struct{}, len(input.Items))
	for _, item := range input.Items {
		if _, err := uuid.Parse(item.ID); err != nil || item.Quantity <= 0 {
			tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Cada artículo debe tener un ID válido y una cantidad mayor que cero."})
			return
		}
		if _, exists := seen[item.ID]; exists {
			tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "No repitas artículos en la solicitud de cierre."})
			return
		}
		seen[item.ID] = struct{}{}
	}

	alreadyProcessed := false
	err := db.GDB.Transaction(func(tx *gorm.DB) error {
		closure := models.InventoryClosure{ID: input.CheckoutID}
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&closure)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			alreadyProcessed = true
			return nil
		}

		for _, item := range input.Items {
			result := tx.Model(&models.InventoryItem{}).
				Where("id = ? AND quantity >= ?", item.ID, item.Quantity).
				UpdateColumn("quantity", gorm.Expr("quantity - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errInsufficientInventory
			}
		}
		return nil
	})
	if errors.Is(err, errInsufficientInventory) {
		tools.JsonResponse(w, http.StatusConflict, map[string]string{"error": "No hay existencias suficientes para uno o más productos. Revisa el inventario antes de cerrar."})
		return
	}
	if err != nil {
		log.Printf("No se pudo descontar el inventario al cerrar %s: %v", input.CheckoutID, err)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo descontar el inventario; el cierre no se completó."})
		return
	}

	tools.JsonResponse(w, http.StatusOK, map[string]bool{"already_processed": alreadyProcessed})
}

var errInsufficientInventory = errors.New("existencias insuficientes")

func GetInventoryItemByID(w http.ResponseWriter, r *http.Request) {
	item, ok := findInventoryItem(w, mux.Vars(r)["id"])
	if ok {
		tools.EncodeJSON(w, toInventoryItemDTO(item))
	}
}

func UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, err := uuid.Parse(id); err != nil {
		tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El ID del artículo no es válido."})
		return
	}

	var input types.InventoryItemInputDTO
	if !tools.DecodeJSON(w, r, &input) {
		return
	}
	validator := validators.NewInventoryValidator()
	if !tools.ValidateAndRespond(w, validator.ValidateItem(&input), validator) {
		return
	}

	result := db.GDB.Model(&models.InventoryItem{}).Where("id = ?", id).Updates(map[string]any{
		"name":       input.Name,
		"quantity":   input.Quantity,
		"net_price":  input.NetPrice,
		"sale_price": input.SalePrice,
	})
	if result.Error != nil {
		log.Printf("No se pudo actualizar el artículo %s: %v", id, result.Error)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo actualizar el artículo."})
		return
	}
	if result.RowsAffected == 0 {
		tools.JsonResponse(w, http.StatusNotFound, map[string]string{"error": "No se encontró el artículo."})
		return
	}

	tools.EncodeJSON(w, types.InventoryItemDTO{
		ID:        id,
		Name:      input.Name,
		Quantity:  input.Quantity,
		NetPrice:  input.NetPrice,
		SalePrice: input.SalePrice,
	})
}

func DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, err := uuid.Parse(id); err != nil {
		tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El ID del artículo no es válido."})
		return
	}

	result := db.GDB.Delete(&models.InventoryItem{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("No se pudo eliminar el artículo %s: %v", id, result.Error)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo eliminar el artículo."})
		return
	}
	if result.RowsAffected == 0 {
		tools.JsonResponse(w, http.StatusNotFound, map[string]string{"error": "No se encontró el artículo."})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func findInventoryItem(w http.ResponseWriter, id string) (models.InventoryItem, bool) {
	if _, err := uuid.Parse(id); err != nil {
		tools.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El ID del artículo no es válido."})
		return models.InventoryItem{}, false
	}

	var item models.InventoryItem
	if err := db.GDB.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tools.JsonResponse(w, http.StatusNotFound, map[string]string{"error": "No se encontró el artículo."})
			return models.InventoryItem{}, false
		}
		log.Printf("No se pudo consultar el artículo %s: %v", id, err)
		tools.JsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "No se pudo cargar el artículo."})
		return models.InventoryItem{}, false
	}
	return item, true
}

func toInventoryItemDTO(item models.InventoryItem) types.InventoryItemDTO {
	return types.InventoryItemDTO{
		ID:        item.ID,
		Name:      item.Name,
		Quantity:  item.Quantity,
		NetPrice:  item.NetPrice,
		SalePrice: item.SalePrice,
	}
}
