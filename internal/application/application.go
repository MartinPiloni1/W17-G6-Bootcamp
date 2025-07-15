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

// SellerRouter creates a new router for seller-related endpoints
func SellerRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()
	repository := repository.NewSellerRepository(db)
	service := service.NewSellerService(repository)
	handler := handler.NewSellerHandler(service)
	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetByID())
	router.Post("/", handler.Create())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())
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

func ProductRouter() chi.Router {
	router := chi.NewRouter()

	rp := repository.NewProductRepositoryFile()
	sv := service.NewProductServiceDefault(rp)
	hd := handler.NewProductHandler(sv)

	router.Post("/", hd.Create())
	router.Get("/", hd.GetAll())
	router.Get("/{id}", hd.GetById())
	router.Patch("/{id}", hd.Update())
	router.Delete("/{id}", hd.Delete())
	return router
}

func BuyersRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	repository := repository.NewBuyerRepositoryDB(db)
	service := service.NewBuyerServiceDefault(repository)
	handler := handler.NewBuyerHandler(service)

	router.Post("/", handler.Create())
	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetByID())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())
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
