package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/application"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type ServerChi struct {
	ServerAddr     string
	DatabaseConfig mysql.Config
}

func LoadServerConf() (*ServerChi, error) {
	// load env file but can be preloaded in the terminal
	err := godotenv.Load()
	if err != nil {
		// just print the error but continue with the execution with docker
		fmt.Println("failed to load env file:", err, " .Assuming preload of env already done")
	}

	// default values
	serverAddr := os.Getenv("ADDRESS")

	if serverAddr == "" {
		serverAddr = "8080"
	}

	Host := os.Getenv("DB_HOST")
	Port := os.Getenv("DB_PORT")
	Name := os.Getenv("DB_NAME")
	User := os.Getenv("DB_USER")
	Pass := os.Getenv("DB_PASS")
	if Host == "" ||
		Port == "" ||
		Name == "" ||
		User == "" {
		return nil, fmt.Errorf("DB conn settings not established")
	}
	dbConfig := storage.NewMySQLConfig(Host, Port, User, Pass, Name)
	serverAddr = ":" + serverAddr // add two point to the address

	return &ServerChi{
		ServerAddr:     serverAddr,
		DatabaseConfig: dbConfig,
	}, nil
}

func (a *ServerChi) Run() (err error) {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // logger

	freshDB, err := storage.InitMySQLConnection(a.DatabaseConfig)
	if err != nil {
		return err
	}

	healthRouter := application.HealthRouter()
	productRouter := application.ProductRouter(freshDB)
	productRecordRouter := application.ProductRecordRouter(freshDB)
	buyersRouter := application.BuyersRouter(freshDB)
	warehouseRouter := application.WarehouseRouter(freshDB)
	sellerRouter := application.SellerRouter(freshDB)
	employeeRouter := application.EmployeeRouter(freshDB)
	sectionRouter := application.SectionRouter(freshDB)
	carryRouter := application.CarryRouter(freshDB)
	localitiesRouter := application.LocalityRouter(freshDB)
	purchaseOrderRouter := application.PurchaseOrderRouter(freshDB)

	router.Mount("/healthcheck", healthRouter)
	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/products", productRouter)
		r.Mount("/productRecords", productRecordRouter)
		r.Mount("/warehouses", warehouseRouter)
		r.Mount("/buyers", buyersRouter)
		r.Mount("/sellers", sellerRouter)
		r.Mount("/employees", employeeRouter)
		r.Mount("/sections", sectionRouter)
		r.Mount("/localities", localitiesRouter)
		r.Mount("/carries", carryRouter)
		r.Mount("/purchaseOrders", purchaseOrderRouter)
	})

	err = http.ListenAndServe(a.ServerAddr, router)
	return
}
