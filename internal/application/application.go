package application

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"net/http"

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

func EmployeeRouter() chi.Router {
	router := chi.NewRouter()

	repo := repository.NewEmployeeRepository()
	service := service.NewEmployeeService(repo)
	handler := handler.NewEmployeeHandler(service)

	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetById())
	router.Post("/", handler.Create())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())
	return router
}
