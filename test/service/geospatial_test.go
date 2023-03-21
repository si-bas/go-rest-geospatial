package test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/si-bas/go-rest-geospatial/domain/model"
	repoMocks "github.com/si-bas/go-rest-geospatial/domain/repository/mocks"
	"github.com/si-bas/go-rest-geospatial/service"
	"github.com/si-bas/go-rest-geospatial/shared/helper/pagination"
	"github.com/stretchr/testify/mock"
	"github.com/twpayne/go-geom/encoding/geojson"
	"gorm.io/gorm"
)

type geospatialMock struct {
	geospatialRepo repoMocks.GeospatialRepository
}

func TestGeospatialList(t *testing.T) {
	geospatials := []model.Geospatial{
		{
			ID:   1,
			Name: "Jakarta",
		},
		{
			ID:   2,
			Name: "Banten",
		},
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *geospatialMock)
		wantErr  error
	}{
		{
			name: "happy flow",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("Get", mock.Anything, model.GeospatialFilter{}).Return(geospatials, nil)
			},
		},
		{
			name: "error - error get geospatial from repo",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("Get", mock.Anything, model.GeospatialFilter{}).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := geospatialMock{
				geospatialRepo: repoMocks.GeospatialRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewGeospatialService(&listMock.geospatialRepo)
			result, err := svc.List(context.TODO(), model.GeospatialFilter{})

			assert.Equal(t, tc.wantErr, err)
			listMock.geospatialRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, geospatials)
			}
		})
	}
}

func TestGeospatialListPaginate(t *testing.T) {
	geospatials := []model.Geospatial{
		{
			ID:   1,
			Name: "Jakarta",
		},
		{
			ID:   2,
			Name: "Banten",
		},
	}

	testCases := []struct {
		name     string
		mockFunc func(mock *geospatialMock)
		wantErr  error
	}{
		{
			name: "happy flow",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetPaginate", mock.Anything, model.GeospatialFilter{}, pagination.Param{}).Return(geospatials, &pagination.Param{}, nil)
			},
		},
		{
			name: "error - error get geospatial from repo",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetPaginate", mock.Anything, model.GeospatialFilter{}, pagination.Param{}).Return(nil, nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := geospatialMock{
				geospatialRepo: repoMocks.GeospatialRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewGeospatialService(&listMock.geospatialRepo)
			result, _, err := svc.ListPaginate(context.TODO(), model.GeospatialFilter{}, pagination.Param{})

			assert.Equal(t, tc.wantErr, err)
			listMock.geospatialRepo.AssertExpectations(t)

			if err == nil {
				assert.Equal(t, result, geospatials)
			}
		})
	}
}

func TestGeospatialGetTypes(t *testing.T) {
	testCases := []struct {
		name     string
		mockFunc func(mock *geospatialMock)
		wantErr  error
	}{
		{
			name: "error - error get types from repo",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetTypes", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "happy flow - data is empty",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetTypes", mock.Anything).Return(nil, nil)
			},
		},
		{
			name: "happy flow",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetTypes", mock.Anything).Return([]string{"A", "B", "C"}, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := geospatialMock{
				geospatialRepo: repoMocks.GeospatialRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewGeospatialService(&listMock.geospatialRepo)
			_, err := svc.GetTypes(context.TODO())

			assert.Equal(t, tc.wantErr, err)
			listMock.geospatialRepo.AssertExpectations(t)
		})
	}
}

func TestGeospatialGetLevels(t *testing.T) {
	testCases := []struct {
		name     string
		mockFunc func(mock *geospatialMock)
		wantErr  error
	}{
		{
			name: "error - error get types from repo",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetLevels", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "happy flow - data is empty",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetLevels", mock.Anything).Return(nil, nil)
			},
		},
		{
			name: "happy flow",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("GetLevels", mock.Anything).Return([]uint{1, 2, 3}, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := geospatialMock{
				geospatialRepo: repoMocks.GeospatialRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewGeospatialService(&listMock.geospatialRepo)
			_, err := svc.GetLevels(context.TODO())

			assert.Equal(t, tc.wantErr, err)
			listMock.geospatialRepo.AssertExpectations(t)
		})
	}
}

func TestCreateFromFeatureCollection(t *testing.T) {
	testCases := []struct {
		name     string
		mockFunc func(mock *geospatialMock)
		wantErr  error
	}{
		{
			name: "Happy flow",
			mockFunc: func(listMock *geospatialMock) {
				listMock.geospatialRepo.On("UpsertBulk", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	// Sample GADM GeoJSON string
	gadmGeoJSON := `{"type":"FeatureCollection","name":"gadm41_IDN_0","crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC:1.3:CRS84"}},"features":[{"type":"Feature","properties":{"GID_0":"IDN","COUNTRY":"Indonesia"},"geometry":{"type":"MultiPolygon","coordinates":[[[[-155.5421143,19.0834808],[-155.6881561,18.9161911],[-155.9368896,19.0593891],[-155.8619995,19.326109]]]]}}]}`

	var features geojson.FeatureCollection
	json.Unmarshal([]byte(gadmGeoJSON), &features)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listMock := geospatialMock{
				geospatialRepo: repoMocks.GeospatialRepository{},
			}
			if tc.mockFunc != nil {
				tc.mockFunc(&listMock)
			}

			svc := service.NewGeospatialService(&listMock.geospatialRepo)
			err := svc.CreateFromFeatureCollection(context.TODO(), &features)

			assert.Equal(t, tc.wantErr, err)
			listMock.geospatialRepo.AssertExpectations(t)
		})
	}
}

func TestGeospatialBuildTree(t *testing.T) {
	geos := []model.Geospatial{
		{ID: 1, GadmID: "1", ParentGadmID: "", Name: "A", Type: "Country", Level: 1},
		{ID: 2, GadmID: "1.1", ParentGadmID: "1", Name: "B", Type: "Province", Level: 2},
		{ID: 3, GadmID: "1.1.1", ParentGadmID: "1.1", Name: "C", Type: "City", Level: 3},
		{ID: 4, GadmID: "1.1.1.1", ParentGadmID: "1.1.1", Name: "D", Type: "District", Level: 4},
	}

	var rootLevel uint = ^uint(0)
	for _, g := range geos {
		if g.Level < rootLevel {
			rootLevel = g.Level
		}
	}

	expectedGeo := &model.Geospatial{
		ID: 1, GadmID: "1", ParentGadmID: "", Name: "A", Type: "Country", Level: 1,
		Child: &model.Geospatial{ID: 2, GadmID: "1.1", ParentGadmID: "1", Name: "B", Type: "Province", Level: 2,
			Child: &model.Geospatial{ID: 3, GadmID: "1.1.1", ParentGadmID: "1.1", Name: "C", Type: "City", Level: 3,
				Child: &model.Geospatial{ID: 4, GadmID: "1.1.1.1", ParentGadmID: "1.1.1", Name: "D", Type: "District", Level: 4}},
		},
	}

	listMock := geospatialMock{
		geospatialRepo: repoMocks.GeospatialRepository{},
	}
	svc := service.NewGeospatialService(&listMock.geospatialRepo)
	result := svc.BuildTree(geos, rootLevel)

	if !reflect.DeepEqual(result, expectedGeo) {
		t.Errorf("Expected %v but got %v", expectedGeo, result)
	}
}
