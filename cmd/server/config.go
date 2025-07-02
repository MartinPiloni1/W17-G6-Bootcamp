package server

import (
	"fmt"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type ServerChi struct {
	ServerAddr     string
	LoaderFilePath string // here whe can add diverse filepaths to the
}

func LoadServerConf() (*ServerChi, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %s", err)
	}
	// default values
	serverAddr := os.Getenv("ADDRESS")

	// here we can load more files if we use it
	filePathDefault := os.Getenv("FILE_PATH_DEFAULT")

	if serverAddr == "" {
		serverAddr = ":8080"
	}

	// here we should validate if they are setted
	if filePathDefault == "" {
		return &ServerChi{}, fmt.Errorf("env variables not setted")
	}

	return &ServerChi{
		ServerAddr:     serverAddr,
		LoaderFilePath: filePathDefault,
	}, nil
}

func (a *ServerChi) Run() (err error) {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // logger

	healthRouter := application.HealthRouter()
	warehouseRouter := application.WarehouseRouter()
	buyersRouter := application.BuyersRouter()

	// mount healthcheck
	router.Mount("/healthcheck", healthRouter)
	router.Mount("/api/v1/warehouses", warehouseRouter)
	router.Mount("/api/v1/buyers", buyersRouter)
	err = http.ListenAndServe(a.ServerAddr, router)
	return
}
