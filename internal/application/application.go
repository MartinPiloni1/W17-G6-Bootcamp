package application

import (
	"net/http"

	"database/sql"
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
	repository := repository.NewSellerRepository()
	service := service.NewSellerService(repository)
	handler := handler.NewSellerHandler(service)
	router.Get("/", handler.GetAll())        // GET /seller         (lista todos)
	router.Get("/{id}", handler.GetByID())   // GET /seller/{id}    (uno por id)
	router.Post("/", handler.Create())       // POST /seller        (crear uno)
	router.Patch("/{id}", handler.Update())  // PATCH /seller/{id}  (actualizar)
	router.Delete("/{id}", handler.Delete()) // DELETE
	return router
}

func WarehouseRouter() chi.Router {
	repository := repository.NewWarehouseRepository()
	service := service.NewWarehouseService(repository)
	handler := handler.NewWarehouseHandler(service)

	router := chi.NewRouter()

	router.Get("/", handler.GetAll())
	router.Post("/", handler.Create())
	router.Get("/{id}", handler.GetById())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())

	return router
}

func ProductRouter() chi.Router {
	router := chi.NewRouter()

	repository := repository.NewProductRepositoryFile()
	service := service.NewProductServiceDefault(repository)
	handler := handler.NewProductHandler(service)

	router.Post("/", handler.Create())
	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetById())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())
	return router
}

func BuyersRouter() chi.Router {
	router := chi.NewRouter()

	repository := repository.NewBuyerRepositoryFile() // fileRepository
	service := service.NewBuyerServiceDefault(repository)
	handler := handler.NewBuyerHandler(service)

	router.Post("/", handler.Create())
	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetByID())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())
	return router
}

func EmployeeRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	repository := repository.NewEmployeeRepository(db)
	service := service.NewEmployeeService(repository)
	handler := handler.NewEmployeeHandler(service)

	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetById())
	router.Post("/", handler.Create())
	router.Patch("/{id}", handler.Update())
	router.Delete("/{id}", handler.Delete())
	return router
}

func SectionRouter() chi.Router {
	repository := repository.NewSectionRepository()
	service := service.NewSectionService(repository)
	handler := handler.NewSectionHandler(service)

	router := chi.NewRouter()

	router.Get("/", handler.GetAll())
	router.Get("/{id}", handler.GetByID())
	router.Delete("/{id}", handler.Delete())
	router.Post("/", handler.Create())
	router.Patch("/{id}", handler.Update())
	return router
}
