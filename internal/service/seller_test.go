package service_test

import (
	"errors"
	"testing"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
)

// TestSellerService_Create tests the Create method of the SellerService.
// It checks for successful creation, conflict errors, and unprocessable entity errors.
func TestSellerService_Create(t *testing.T) {
	attrOK := models.SellerAttributes{CID: 1, CompanyName: "A", Address: "B", Telephone: "C", LocalityID: "L"}
	attrConflict := models.SellerAttributes{CID: 2, CompanyName: "B", Address: "Z", Telephone: "W", LocalityID: "L"}

	cases := []struct {
		name        string
		input       models.SellerAttributes
		setupMock   func(r *mocks.SellerRepositoryDBMock)
		wantResult  models.Seller
		wantErr     bool
		errContains string
	}{
		{
			name:  "create_ok",
			input: attrOK,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("Create", attrOK).Return(models.Seller{ID: 99, SellerAttributes: attrOK}, nil)
			},
			wantResult:  models.Seller{ID: 99, SellerAttributes: attrOK},
			wantErr:     false,
			errContains: "",
		},
		{
			name:  "create_conflict",
			input: attrConflict,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("Create", attrConflict).Return(models.Seller{}, httperrors.ConflictError{Message: "CID already exists"})
			},
			wantResult:  models.Seller{},
			wantErr:     true,
			errContains: "CID already exists",
		},
		{
			name:        "unprocessable_entity",
			input:       models.SellerAttributes{},
			setupMock:   func(r *mocks.SellerRepositoryDBMock) {},
			wantResult:  models.Seller{},
			wantErr:     true,
			errContains: "Invalid seller data",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.SellerRepositoryDBMock)
			tc.setupMock(repository)
			service := service.NewSellerService(repository)
			result, err := service.Create(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				assert.Equal(t, tc.wantResult, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, result)
			}
			repository.AssertExpectations(t)
		})
	}
}

// TestSellerService_GetAll tests the GetAll method of the SellerService.
// It checks for successful retrieval of all sellers and error handling.
func TestSellerService_GetAll(t *testing.T) {
	list := []models.Seller{
		{ID: 1, SellerAttributes: models.SellerAttributes{CID: 10}},
		{ID: 2, SellerAttributes: models.SellerAttributes{CID: 20}},
	}
	cases := []struct {
		name      string
		setupMock func(r *mocks.SellerRepositoryDBMock)
		wantRes   []models.Seller
		wantErr   bool
	}{
		{
			name: "find_all",
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("GetAll").Return(list, nil)
			},
			wantRes: list,
			wantErr: false,
		},
		{
			name: "find_all_fail",
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("GetAll").Return([]models.Seller{}, errors.New("db error"))
			},
			wantRes: []models.Seller{},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.SellerRepositoryDBMock)
			tc.setupMock(repository)
			service := service.NewSellerService(repository)
			got, err := service.GetAll()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantRes, got)
			repository.AssertExpectations(t)
		})
	}
}

// TestSellerService_GetByID tests the GetByID method of the SellerService.
// It checks for successful retrieval of a seller by ID and error handling for non-existent sellers.
func TestSellerService_GetByID(t *testing.T) {
	cases := []struct {
		name        string
		input       int
		setupMock   func(r *mocks.SellerRepositoryDBMock)
		wantResult  models.Seller
		wantErr     bool
		errContains string
	}{
		{
			name:  "find_by_id_existent",
			input: 5,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				s := models.Seller{ID: 5, SellerAttributes: models.SellerAttributes{CID: 55}}
				r.On("GetByID", 5).Return(s, nil)
			},
			wantResult: models.Seller{ID: 5, SellerAttributes: models.SellerAttributes{CID: 55}},
			wantErr:    false,
		},
		{
			name:  "find_by_id_non_existent",
			input: 99,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("GetByID", 99).Return(models.Seller{}, httperrors.NotFoundError{Message: "not found"})
			},
			wantResult:  models.Seller{},
			wantErr:     true,
			errContains: "not found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(mocks.SellerRepositoryDBMock)
			tc.setupMock(repository)
			service := service.NewSellerService(repository)
			res, err := service.GetByID(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				assert.Equal(t, tc.wantResult, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, res)
			}
			repository.AssertExpectations(t)
		})
	}
}

// TestSellerService_Delete tests the Delete method of the SellerService.
// It checks for successful deletion of a seller and error handling for non-existent sellers.
func TestSellerService_Delete(t *testing.T) {
	cases := []struct {
		name        string
		id          int
		setupMock   func(r *mocks.SellerRepositoryDBMock)
		wantErr     bool
		errContains string
	}{
		{
			name: "delete_ok",
			id:   1,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("Delete", 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete_non_existent",
			id:   5,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("Delete", 5).Return(httperrors.NotFoundError{Message: "Seller not found"})
			},
			wantErr:     true,
			errContains: "Seller not found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := new(mocks.SellerRepositoryDBMock)
			tc.setupMock(r)
			svc := service.NewSellerService(r)
			err := svc.Delete(tc.id)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			r.AssertExpectations(t)
		})
	}
}

// TestSellerService_Update tests the Update method of the SellerService.
// It checks for successful updates, conflict errors, and non-existent sellers.
func TestSellerService_Update(t *testing.T) {
	attr := &models.SellerAttributes{CID: 42, CompanyName: "UPD", Address: "A", Telephone: "32", LocalityID: "X"}
	updated := models.Seller{ID: 2, SellerAttributes: *attr}

	cases := []struct {
		name        string
		id          int
		input       *models.SellerAttributes
		setupMock   func(r *mocks.SellerRepositoryDBMock)
		wantResult  models.Seller
		wantErr     bool
		errContains string
	}{
		{
			name:  "update_ok",
			id:    2,
			input: attr,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("GetAll").Return([]models.Seller{}, nil)
				r.On("Update", 2, attr).Return(updated, nil)
			},
			wantResult: updated,
			wantErr:    false,
		},
		{
			name:  "update_non_existent",
			id:    77,
			input: attr,
			setupMock: func(r *mocks.SellerRepositoryDBMock) {
				r.On("GetAll").Return([]models.Seller{}, nil)
				r.On("Update", 77, attr).Return(models.Seller{}, httperrors.NotFoundError{Message: "Seller not found"})
			},
			wantResult:  models.Seller{},
			wantErr:     true,
			errContains: "Seller not found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := new(mocks.SellerRepositoryDBMock)
			tc.setupMock(r)
			svc := service.NewSellerService(r)
			res, err := svc.Update(tc.id, tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				assert.Equal(t, tc.wantResult, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, res)
			}
			r.AssertExpectations(t)
		})
	}
}
