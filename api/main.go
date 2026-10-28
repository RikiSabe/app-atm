package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"app-atm-api/internal/db"
	"app-atm-api/internal/routers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("No se pudo cargar .env: %v", err)
	}

	if err := db.Connection(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	if err := db.MigrateInventory(); err != nil {
		log.Fatalf("Error al migrar los modelos de la base de datos: %v", err)
	}

	router := mux.NewRouter()
	routers.InitEndPoints(router)

	allowedOrigins := []string{
		"http://localhost:3000",
		"https://app-atm.vercel.app",
	}
	if configuredOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); configuredOrigins != "" {
		allowedOrigins = strings.Split(configuredOrigins, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	handler := handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedHeaders([]string{"Content-Type", "Accept", "Origin"}),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		}),
	)(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5665"
	}
	address := "0.0.0.0:" + port
	log.Printf("API de inventario escuchando en %s", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
