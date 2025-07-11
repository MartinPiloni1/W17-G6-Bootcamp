package server

import (
	config "github.com/aaguero_meli/W17-G6-Bootcamp/cmd/db"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerChi struct {
	ServerAddr     string
	DatabaseConfig mysql.Config
}

func LoadServerConf() (*ServerChi, error) {
	config.LoadEnv() // Cargar .env sólo una vez

	// La dirección, sigue igual:
	serverAddr := os.Getenv("ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	// Delegar la parte de DB en la función común:
	dbConfig := config.GetDBConfigFromEnv()

	return &ServerChi{
		ServerAddr:     serverAddr,
		DatabaseConfig: dbConfig,
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
