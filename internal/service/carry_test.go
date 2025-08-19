package service_test

import (
	"errors"
	"testing"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getValidCarryAttributes() models.CarryAttributes {
	return models.CarryAttributes{
		Cid:         "CARRY001",
		CompanyName: "Transportes Rápidos S.A.",
		Address:     "Av. Principal 123, Ciudad",
		Telephone:   "+54-11-1234-5678",
		LocalityId:  "LOC001",
	}
}

// TestCarryServiceDefault_Create tests the Create method of the CarryService.
// It covers successful creation, validation errors, and repository errors.
func TestCarryServiceDefault_Create(t *testing.T) {
	validAttrs := getValidCarryAttributes()
	expectedCarry := models.Carry{
		Id:              1,
		CarryAttributes: validAttrs,
	}

	t.Run("create_ok: success to create carry", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		mockRepo.On("Create", validAttrs).Return(expectedCarry, nil)

		result, err := carryService.Create(validAttrs)

		require.NoError(t, err)
		assert.Equal(t, expectedCarry, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_empty_cid: returns bad request error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		invalidAttrs := validAttrs
		invalidAttrs.Cid = ""

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Cid must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_empty_company_name: returns bad request error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		invalidAttrs := validAttrs
		invalidAttrs.CompanyName = ""

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field CompanyName must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_empty_address: returns bad request error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		invalidAttrs := validAttrs
		invalidAttrs.Address = ""

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Address must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_empty_telephone: returns bad request error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		invalidAttrs := validAttrs
		invalidAttrs.Telephone = ""

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Telephone must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_empty_locality_id: returns bad request error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		invalidAttrs := validAttrs
		invalidAttrs.LocalityId = ""

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field LocalityId must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_multiple_empty_fields: returns first validation error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		invalidAttrs := models.CarryAttributes{
			Cid:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityId:  "",
		}

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Cid must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_repository_conflict_error: returns conflict error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		conflictError := httperrors.ConflictError{Message: "the Cid already exists"}
		mockRepo.On("Create", validAttrs).Return(models.Carry{}, conflictError)

		result, err := carryService.Create(validAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.ConflictError{})
		assert.Contains(t, err.Error(), "the Cid already exists")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_repository_locality_not_found: returns conflict error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		localityError := httperrors.ConflictError{Message: "the LocalityId does not exist"}
		mockRepo.On("Create", validAttrs).Return(models.Carry{}, localityError)

		result, err := carryService.Create(validAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.ConflictError{})
		assert.Contains(t, err.Error(), "the LocalityId does not exist")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_repository_internal_error: returns internal server error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		internalError := httperrors.InternalServerError{Message: "error creating carry"}
		mockRepo.On("Create", validAttrs).Return(models.Carry{}, internalError)

		result, err := carryService.Create(validAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Contains(t, err.Error(), "error creating carry")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_repository_generic_error: returns generic error", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		genericError := errors.New("database connection failed")
		mockRepo.On("Create", validAttrs).Return(models.Carry{}, genericError)

		result, err := carryService.Create(validAttrs)

		require.Error(t, err)
		assert.Equal(t, genericError, err)
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})
}

// TestCarryServiceDefault_Create_EdgeCases tests edge cases for carry creation
func TestCarryServiceDefault_Create_EdgeCases(t *testing.T) {
	t.Run("create_with_whitespace_only_fields: should pass validation", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		attrsWithWhitespace := models.CarryAttributes{
			Cid:         "   CARRY002   ",
			CompanyName: "   Company with spaces   ",
			Address:     "   Address with spaces   ",
			Telephone:   "   +54-11-1234-5678   ",
			LocalityId:  "   LOC002   ",
		}

		expectedCarry := models.Carry{
			Id:              2,
			CarryAttributes: attrsWithWhitespace,
		}

		mockRepo.On("Create", attrsWithWhitespace).Return(expectedCarry, nil)

		result, err := carryService.Create(attrsWithWhitespace)

		require.NoError(t, err)
		assert.Equal(t, expectedCarry, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_with_special_characters: should pass validation", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		attrsWithSpecialChars := models.CarryAttributes{
			Cid:         "CARRY-003@#$%",
			CompanyName: "Empresa & Asociados S.A.",
			Address:     "Calle #123, Piso 4°, Depto. 5",
			Telephone:   "+54-11-1234-5678 ext. 123",
			LocalityId:  "LOC-003_ABC",
		}

		expectedCarry := models.Carry{
			Id:              3,
			CarryAttributes: attrsWithSpecialChars,
		}

		mockRepo.On("Create", attrsWithSpecialChars).Return(expectedCarry, nil)

		result, err := carryService.Create(attrsWithSpecialChars)

		require.NoError(t, err)
		assert.Equal(t, expectedCarry, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_with_long_values: should pass validation", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		longCompanyName := "Empresa de Transporte Internacional con Nombre Muy Largo y Complejo para Pruebas de Validación de Longitud Máxima de Campo en el Sistema de Gestión de Carriers"
		longAddress := "Avenida Principal de la Ciudad de Buenos Aires, Número 12345, Piso 15, Departamento 1501, Edificio Torre Empresarial, Zona Comercial Norte, Código Postal 1425, República Argentina"

		attrsWithLongValues := models.CarryAttributes{
			Cid:         "CARRY004",
			CompanyName: longCompanyName,
			Address:     longAddress,
			Telephone:   "+54-11-1234-5678",
			LocalityId:  "LOC004",
		}

		expectedCarry := models.Carry{
			Id:              4,
			CarryAttributes: attrsWithLongValues,
		}

		mockRepo.On("Create", attrsWithLongValues).Return(expectedCarry, nil)

		result, err := carryService.Create(attrsWithLongValues)

		require.NoError(t, err)
		assert.Equal(t, expectedCarry, result)
		mockRepo.AssertExpectations(t)
	})
}

// TestCarryServiceDefault_Create_ValidationOrder tests that validation errors are returned in the correct order
func TestCarryServiceDefault_Create_ValidationOrder(t *testing.T) {
	t.Run("validation_order_cid_first: should return cid error first", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		// All fields empty, should return Cid error first
		invalidAttrs := models.CarryAttributes{
			Cid:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityId:  "",
		}

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Cid must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation_order_company_name_second: should return company name error second", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		// Cid valid, others empty, should return CompanyName error
		invalidAttrs := models.CarryAttributes{
			Cid:         "CARRY005",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityId:  "",
		}

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field CompanyName must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation_order_address_third: should return address error third", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		// Cid and CompanyName valid, others empty, should return Address error
		invalidAttrs := models.CarryAttributes{
			Cid:         "CARRY006",
			CompanyName: "Valid Company",
			Address:     "",
			Telephone:   "",
			LocalityId:  "",
		}

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Address must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation_order_telephone_fourth: should return telephone error fourth", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		// Cid, CompanyName, and Address valid, others empty, should return Telephone error
		invalidAttrs := models.CarryAttributes{
			Cid:         "CARRY007",
			CompanyName: "Valid Company",
			Address:     "Valid Address",
			Telephone:   "",
			LocalityId:  "",
		}

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field Telephone must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation_order_locality_id_last: should return locality id error last", func(t *testing.T) {
		mockRepo := new(mocks.CarryRepositoryDBMock)
		carryService := service.NewCarryService(mockRepo)

		// All fields valid except LocalityId, should return LocalityId error
		invalidAttrs := models.CarryAttributes{
			Cid:         "CARRY008",
			CompanyName: "Valid Company",
			Address:     "Valid Address",
			Telephone:   "+54-11-1234-5678",
			LocalityId:  "",
		}

		result, err := carryService.Create(invalidAttrs)

		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.BadRequestError{})
		assert.Contains(t, err.Error(), "the field LocalityId must not be empty")
		assert.Equal(t, models.Carry{}, result)
		mockRepo.AssertExpectations(t)
	})
}
