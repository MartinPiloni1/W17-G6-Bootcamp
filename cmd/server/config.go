package server

import (
	"log"
	"net/http"
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type ServerChi struct {
	ServerAddr string
}

func LoadServerConf() (*ServerChi, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %s", err)
	}
	// default values
	serverAddr := os.Getenv("ADDRESS")

	if serverAddr == "" {
		serverAddr = ":8080"
	}

	return &ServerChi{
		ServerAddr: serverAddr,
	}, nil
}

func (a *ServerChi) Run() (err error) {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // logger

	healthRouter := application.HealthRouter()
	productRouter := application.ProductRouter()
	warehouseRouter := application.WarehouseRouter()
	buyersRouter := application.BuyersRouter()
	sellerRouter := application.SellerRouter()
	employeeRouter := application.EmployeeRouter()
	sectionRouter := application.SectionRouter()

	router.Mount("/healthcheck", healthRouter)
	router.Mount("/api/v1/products", productRouter)
	router.Mount("/api/v1/warehouses", warehouseRouter)
	router.Mount("/api/v1/buyers", buyersRouter)
	router.Mount("/api/v1/sellers", sellerRouter)
	router.Mount("/api/v1/employees", employeeRouter)
	router.Mount("/api/v1/sections", sectionRouter)
	err = http.ListenAndServe(a.ServerAddr, router)
	return
}
