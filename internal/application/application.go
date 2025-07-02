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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	})
	return router
}

func BuyersRouter() chi.Router {
	router := chi.NewRouter()

	rp := repository.NewBuyerRepositoryFile() // fileRepository
	sv := service.NewBuyerServiceDefault(rp)
	hd := handler.NewBuyerHandler(sv)

	router.Post("/", hd.Create())
	router.Get("/", hd.GetAll())
	router.Get("/{id}", hd.GetByID())
	router.Patch("/{id}", hd.Update())
	router.Delete("/{id}", hd.Delete())
	return router
}
