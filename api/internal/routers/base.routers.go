package routers

import "github.com/gorilla/mux"

func InitEndPoints(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()
	InventoryRoutes(api)
}
