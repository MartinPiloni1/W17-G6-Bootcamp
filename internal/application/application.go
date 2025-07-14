package application

import (
	"database/sql"
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

func WarehouseRouter() chi.Router {
	rp := repository.NewWarehouseRepository()
	sv := service.NewWarehouseService(rp)
	hd := handler.NewWarehouseHandler(sv)

	router := chi.NewRouter()

	router.Get("/", hd.GetAll())
	router.Post("/", hd.Create())
	router.Get("/{id}", hd.GetById())
	router.Patch("/{id}", hd.Update())
	router.Delete("/{id}", hd.Delete())

	return router
}

func ProductRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	rp := repository.NewProductRepositoryDB(db)
	sv := service.NewProductServiceDefault(rp)
	hd := handler.NewProductHandler(sv)

	router.Post("/", hd.Create())
	router.Get("/", hd.GetAll())
	router.Get("/{id}", hd.GetById())
	router.Patch("/{id}", hd.Update())
	router.Delete("/{id}", hd.Delete())
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

func EmployeeRouter() chi.Router {
	router := chi.NewRouter()

	rp := repository.NewEmployeeRepository()
	sv := service.NewEmployeeService(rp)
	hd := handler.NewEmployeeHandler(sv)

	router.Get("/", hd.GetAll())
	router.Get("/{id}", hd.GetById())
	router.Post("/", hd.Create())
	router.Patch("/{id}", hd.Update())
	router.Delete("/{id}", hd.Delete())
	return router
}

func SectionRouter() chi.Router {
	rp := repository.NewSectionRepository()
	sv := service.NewSectionService(rp)
	hd := handler.NewSectionHandler(sv)

	router := chi.NewRouter()

	router.Get("/", hd.GetAll())
	router.Get("/{id}", hd.GetByID())
	router.Delete("/{id}", hd.Delete())
	router.Post("/", hd.Create())
	router.Patch("/{id}", hd.Update())
	return router
}
