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

func TestLocalityService_Create(t *testing.T) {
	valid := models.Locality{ID: "L1", LocalityName: "Loc", ProvinceName: "BA", CountryName: "Argentina"}
	missing := models.Locality{}
	cases := []struct {
		name        string
		input       models.Locality
		setupMock   func(m *mocks.LocalityRepositoryMock)
		wantErr     bool
		errContains string
	}{
		{
			name:  "create_ok",
			input: valid,
			setupMock: func(m *mocks.LocalityRepositoryMock) {
				m.On("Create", valid).Return(valid, nil)
			},
			wantErr: false,
		},
		{
			name:        "create_fail_missing_fields",
			input:       missing,
			setupMock:   func(m *mocks.LocalityRepositoryMock) {},
			wantErr:     true,
			errContains: "Invalid locality data",
		},
		{
			name:  "create_conflict",
			input: valid,
			setupMock: func(m *mocks.LocalityRepositoryMock) {
				m.On("Create", valid).Return(models.Locality{}, httperrors.ConflictError{Message: "ID exists"})
			},
			wantErr:     true,
			errContains: "ID exists",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.LocalityRepositoryMock)
			tc.setupMock(repo)
			svc := service.NewLocalityService(repo)
			_, err := svc.Create(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestLocalityService_GetByID(t *testing.T) {
	ok := models.Locality{ID: "A", LocalityName: "LocA", ProvinceName: "BA", CountryName: "AR"}
	cases := []struct {
		name        string
		id          string
		setupMock   func(m *mocks.LocalityRepositoryMock)
		wantResult  models.Locality
		wantErr     bool
		errContains string
	}{
		{
			name:       "found",
			id:         "A",
			setupMock:  func(m *mocks.LocalityRepositoryMock) { m.On("GetByID", "A").Return(ok, nil) },
			wantResult: ok,
			wantErr:    false,
		},
		{
			name: "not_found",
			id:   "B",
			setupMock: func(m *mocks.LocalityRepositoryMock) {
				m.On("GetByID", "B").Return(models.Locality{}, httperrors.NotFoundError{Message: "not found"})
			},
			wantResult:  models.Locality{},
			wantErr:     true,
			errContains: "not found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.LocalityRepositoryMock)
			tc.setupMock(repo)
			svc := service.NewLocalityService(repo)
			res, err := svc.GetByID(tc.id)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, res)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestLocalityService_GetSellerReport(t *testing.T) {
	report := []models.SellerReport{
		{LocalityID: "A", LocalityName: "LocA", SellersCount: 5},
	}
	cases := []struct {
		name        string
		id          *string
		setupMock   func(m *mocks.LocalityRepositoryMock)
		wantResult  []models.SellerReport
		wantErr     bool
		errContains string
	}{
		{
			name:       "ok",
			id:         nil,
			setupMock:  func(m *mocks.LocalityRepositoryMock) { m.On("GetSellerReport", (*string)(nil)).Return(report, nil) },
			wantResult: report,
			wantErr:    false,
		},
		{
			name: "repo_fail",
			id:   nil,
			setupMock: func(m *mocks.LocalityRepositoryMock) {
				m.On("GetSellerReport", (*string)(nil)).Return(nil, errors.New("db error"))
			},
			wantResult:  nil,
			wantErr:     true,
			errContains: "db error",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.LocalityRepositoryMock)
			tc.setupMock(repo)
			svc := service.NewLocalityService(repo)
			res, err := svc.GetSellerReport(tc.id)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, res)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestLocalityService_GetReportByLocalityId(t *testing.T) {
	carry := []models.CarryReport{
		{LocalityId: "Z", LocalityName: "LocZ", CarriesCount: 1},
	}
	cases := []struct {
		name        string
		id          string
		setupMock   func(m *mocks.LocalityRepositoryMock)
		wantResult  []models.CarryReport
		wantErr     bool
		errContains string
	}{
		{
			name:       "ok",
			id:         "Z",
			setupMock:  func(m *mocks.LocalityRepositoryMock) { m.On("GetReportByLocalityId", "Z").Return(carry, nil) },
			wantResult: carry,
			wantErr:    false,
		},
		{
			name: "not_found",
			id:   "W",
			setupMock: func(m *mocks.LocalityRepositoryMock) {
				m.On("GetReportByLocalityId", "W").Return(nil, nil)
			},
			wantResult:  nil,
			wantErr:     true,
			errContains: "locality not found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.LocalityRepositoryMock)
			tc.setupMock(repo)
			svc := service.NewLocalityService(repo)
			res, err := svc.GetReportByLocalityId(tc.id)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, res)
			}
			repo.AssertExpectations(t)
		})
	}
}
