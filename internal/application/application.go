package application

import (
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/go-chi/chi/v5"
)

func HealthRouter() chi.Router {
	// here we cand add the storage (if any) repository service and handler
	//mount the endpoints and return to mount them in the Run()
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return router
}

func SellerRouter() chi.Router {
	router := chi.NewRouter()
	repo := repository.NewSellerRepository()
	service := service.NewSellerService(repo)
	handler := handler.NewSellerHandler(service)
	router.Get("/getall", handler.GetAll())
	return router
}
