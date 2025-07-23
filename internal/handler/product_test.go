package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mock/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/require"
)

// TODO: Move this function to an utils package
func Ptr[T any](v T) *T { return &v }

func TestProductHandler(t *testing.T) {
	t.Run("GetAll returns every entry of the DB", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		products := []models.Product{
			{
				ID: 1,
				ProductAttributes: models.ProductAttributes{
					Description:                    "Yogurt helado",
					ExpirationRate:                 7,
					FreezingRate:                   2,
					Height:                         10.5,
					Length:                         20.0,
					Width:                          15.0,
					NetWeight:                      1.2,
					ProductCode:                    "YOG01",
					RecommendedFreezingTemperature: -5.0,
					ProductTypeID:                  3,
					SellerID:                       Ptr(1),
				},
			},
			{
				ID: 2,
				ProductAttributes: models.ProductAttributes{
					Description:                    "Pechuga de pollo",
					ExpirationRate:                 3,
					FreezingRate:                   1,
					Height:                         5.0,
					Length:                         25.0,
					Width:                          12.5,
					NetWeight:                      0.8,
					ProductCode:                    "POLLO01",
					RecommendedFreezingTemperature: 0.0,
					ProductTypeID:                  7,
					SellerID:                       Ptr(2),
				},
			},
		}

		serviceMock := &mocks.ProductServiceMock{}
		serviceMock.On("GetAll", ctx).Return(products, nil)

		hd := handler.NewProductHandler(serviceMock)
		getAllFunc := hd.GetAll()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		expectedCode := http.StatusOK
		expectedBody := `
		{
			"data": [
				{
					"id": 1,
					"description": "Yogurt helado",
					"expiration_rate": 7,
					"freezing_rate": 2,
					"height": 10.5,
					"length": 20.0,
					"width": 15.0,
					"netweight": 1.2,
					"product_code": "YOG01",
					"recommended_freezing_temperature": -5.0,
					"product_type_id": 3,
					"seller_id": 1 
				},
				{
					"id": 2,
					"description": "Pechuga de pollo",
					"expiration_rate": 3,
					"freezing_rate": 1,
					"height": 5.0,
					"length": 25.0,
					"width": 12.5,
					"netweight": 0.8,
					"product_code": "POLLO01",
					"recommended_freezing_temperature": 0.0,
					"product_type_id": 7,
					"seller_id": 2  
				}
			]
		}`

		// Act
		getAllFunc(response, request)

		// Assert
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}
