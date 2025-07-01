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
	rp := repository.NewSellerRepository()
	sv := service.NewSellerService(rp)
	hd := handler.NewSellerHandler(sv)
	router.Get("/", hd.GetAll())        // GET /seller         (lista todos)
	router.Get("/{id}", hd.GetByID())   // GET /seller/{id}    (uno por id)
	router.Post("/", hd.Create())       // POST /seller        (crear uno)
	router.Patch("/{id}", hd.Update())  // PATCH /seller/{id}  (actualizar)
	router.Delete("/{id}", hd.Delete()) // DELETE
	return router
}
