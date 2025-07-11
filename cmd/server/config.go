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
	ServerAddr string
	DBConf     mysql.Config
}

func LoadServerConf(withEnvFile bool) (*ServerChi, error) {
	if withEnvFile {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("Failed to load env file: %w", err)
		}
	}
	// default values
	serverAddr := os.Getenv("ADDRESS")

	if serverAddr == "" {
		serverAddr = ":8080"
	}

	Host := os.Getenv("DB_HOST")
	Port := os.Getenv("DB_PORT")
	Name := os.Getenv("DB_NAME")
	User := os.Getenv("DB_USER")
	Pass := os.Getenv("DB_PASS")
	if Host == "" ||
		Port == "" ||
		Name == "" ||
		User == "" ||
		Pass == "" {
		return nil, fmt.Errorf("DB conn settings not established")
	}
	dbConf := storage.NewMySQLConfig(Host, Port, User, Pass, Name)

	return &ServerChi{
		ServerAddr: serverAddr,
		DBConf:     dbConf,
	}, nil
}

func (a *ServerChi) Run() (err error) {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // logger

	freshDB, err := storage.InitMySQLConnection(a.DBConf)
	if err != nil {
		return err
	}
	fmt.Println(freshDB) // TODO: remove when in use

	healthRouter := application.HealthRouter()
	productRouter := application.ProductRouter()
	warehouseRouter := application.WarehouseRouter()
	buyersRouter := application.BuyersRouter()
	sellerRouter := application.SellerRouter()
	employeeRouter := application.EmployeeRouter()
	sectionRouter := application.SectionRouter()

	router.Mount("/healthcheck", healthRouter)
	router.Route("/api/v1", func(r chi.Router) {
		router.Mount("/products", productRouter)
		router.Mount("/warehouses", warehouseRouter)
		router.Mount("/buyers", buyersRouter)
		router.Mount("/sellers", sellerRouter)
		router.Mount("/employees", employeeRouter)
		router.Mount("/sections", sectionRouter)
	})

	err = http.ListenAndServe(a.ServerAddr, router)
	return
}
