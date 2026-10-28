package routers

import (
	"app-atm-api/internal/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func InventoryRoutes(api *mux.Router) {
	inventory := api.PathPrefix("/inventory").Subrouter()
	inventory.HandleFunc("", controllers.GetInventoryItems).Methods(http.MethodGet)
	inventory.HandleFunc("", controllers.CreateInventoryItem).Methods(http.MethodPost)
	inventory.HandleFunc("/for-sale", controllers.GetSellableInventoryItems).Methods(http.MethodGet)
	inventory.HandleFunc("/consume", controllers.ConsumeInventoryItems).Methods(http.MethodPost)
	inventory.HandleFunc("/{id}", controllers.GetInventoryItemByID).Methods(http.MethodGet)
	inventory.HandleFunc("/{id}", controllers.UpdateInventoryItem).Methods(http.MethodPut)
	inventory.HandleFunc("/{id}", controllers.DeleteInventoryItem).Methods(http.MethodDelete)
}
