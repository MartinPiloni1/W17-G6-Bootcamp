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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}

// SellerRouter creates a new router for seller-related endpoints
func SellerRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()
	sellerRepository := repository.NewSellerRepository(db)
	sellerService := service.NewSellerService(sellerRepository)
	sellerHandler := handler.NewSellerHandler(sellerService)
	router.Get("/", sellerHandler.GetAll())
	router.Get("/{id}", sellerHandler.GetByID())
	router.Post("/", sellerHandler.Create())
	router.Patch("/{id}", sellerHandler.Update())
	router.Delete("/{id}", sellerHandler.Delete())
	return router
}

func LocalityRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	localityRepository := repository.NewLocalityRepository(db)
	localityServer := service.NewLocalityService(localityRepository)
	localityHandler := handler.NewLocalityHandler(localityServer)

	router.Post("/", localityHandler.Create())
	router.Get("/{id}", localityHandler.GetByID())
	router.Get("/reportSellers", localityHandler.GetSellerReport())
	router.Get("/reportCarries", localityHandler.GetReportByLocalityId())
	return router
}

func WarehouseRouter(db *sql.DB) chi.Router {
	warehouseRepository := repository.NewWarehouseRepositoryDb(db)
	warehouseService := service.NewWarehouseService(warehouseRepository)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)

	router := chi.NewRouter()

	router.Get("/", warehouseHandler.GetAll())
	router.Post("/", warehouseHandler.Create())
	router.Get("/{id}", warehouseHandler.GetById())
	router.Patch("/{id}", warehouseHandler.Update())
	router.Delete("/{id}", warehouseHandler.Delete())

	return router
}

// ProductRouter creates and returns a chi.Router configured
// with CRUD endpoints for products.
func ProductRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	productRepository := repository.NewProductRepositoryDB(db)
	productService := service.NewProductServiceDefault(productRepository)
	productHandler := handler.NewProductHandler(productService)

	router.Post("/", productHandler.Create())
	router.Get("/", productHandler.GetAll())
	router.Get("/{id}", productHandler.GetById())
	router.Get("/reportRecords", productHandler.GetRecordsPerProduct())
	router.Patch("/{id}", productHandler.Update())
	router.Delete("/{id}", productHandler.Delete())
	return router
}

// ProductRecordRouter creates and returns a chi.Router configured for product_records.
func ProductRecordRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	productRecordRepository := repository.NewProductRecordRepositoryDB(db)
	productRecordService := service.NewProductRecordServiceDefault(productRecordRepository)
	productRecordHandler := handler.NewProductRecordHandler(productRecordService)

	router.Post("/", productRecordHandler.Create())
	return router
}

func BuyersRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	buyersRepository := repository.NewBuyerRepositoryDB(db)
	buyersService := service.NewBuyerServiceDefault(buyersRepository)
	buyersHandler := handler.NewBuyerHandler(buyersService)

	router.Post("/", buyersHandler.Create())
	router.Get("/", buyersHandler.GetAll())
	router.Get("/{id}", buyersHandler.GetByID())
	router.Patch("/{id}", buyersHandler.Update())
	router.Delete("/{id}", buyersHandler.Delete())
	router.Get("/reportPurchaseOrders", buyersHandler.GetWithPurchaseOrdersCount())
	return router
}

func EmployeeRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	employeeRepository := repository.NewEmployeeRepository(db)
	inboundOrderRepository := repository.NewInboundOrderRepository(db)
	employeeService := service.NewEmployeeService(employeeRepository, inboundOrderRepository)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	router.Get("/", employeeHandler.GetAll())
	router.Get("/{id}", employeeHandler.GetById())
	router.Post("/", employeeHandler.Create())
	router.Patch("/{id}", employeeHandler.Update())
	router.Delete("/{id}", employeeHandler.Delete())
	router.Get("/reportInboundOrders", employeeHandler.GetInboundOrderReport())
	return router
}

func SectionRouter(db *sql.DB) chi.Router {
	sectionRepository := repository.NewSectionRepositoryDB(db)
	sectionService := service.NewSectionServiceDefault(sectionRepository)
	sectionHandler := handler.NewSectionHandler(sectionService)

	router := chi.NewRouter()

	router.Get("/", sectionHandler.GetAll())
	router.Get("/{id}", sectionHandler.GetByID())
	router.Delete("/{id}", sectionHandler.Delete())
	router.Post("/", sectionHandler.Create())
	router.Patch("/{id}", sectionHandler.Update())
	router.Get("/reportProducts", sectionHandler.GetProductsReport())
	return router
}

func PurchaseOrderRouter(db *sql.DB) chi.Router {
	purchaseOrderRepository := repository.NewPurchaseOrderRepositoryDB(db)
	purchaseOrderService := service.NewPurchaseOrderDefault(purchaseOrderRepository)
	purchaseOrderHandler := handler.NewPurchaseOrderHandler(purchaseOrderService)

	router := chi.NewRouter()

	router.Post("/", purchaseOrderHandler.Create())
	return router
}

func CarryRouter(db *sql.DB) chi.Router {
	carryRepository := repository.NewCarryRepositoryDb(db)
	carryService := service.NewCarryService(carryRepository)
	carryHandler := handler.NewCarryHandler(carryService)

	router := chi.NewRouter()
	router.Post("/", carryHandler.Create())

	return router
}

func InboundOrderRouter(db *sql.DB) chi.Router {
	router := chi.NewRouter()

	inboundOrderRepository := repository.NewInboundOrderRepository(db)
	employeeRepository := repository.NewEmployeeRepository(db)
	warehouseRepository := repository.NewWarehouseRepositoryDb(db)
	service := service.NewInboundOrderService(inboundOrderRepository, employeeRepository, warehouseRepository)
	handler := handler.NewInboundOrderHandler(service)

	router.Post("/", handler.Create())
	return router
}

func ProductBatchRouter(db *sql.DB) chi.Router {
	productBatchRepository := repository.NewProductBatchRepositoryDB(db)
	productBatchService := service.NewProductBatchServiceDefault(productBatchRepository)
	productBatchHandler := handler.NewProductBatchHandler(productBatchService)

	router := chi.NewRouter()
	router.Post("/", productBatchHandler.Create())
	return router
}
