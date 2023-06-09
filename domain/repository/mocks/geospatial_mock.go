// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "github.com/si-bas/go-rest-geospatial/domain/model"

	pagination "github.com/si-bas/go-rest-geospatial/shared/helper/pagination"
)

// GeospatialRepository is an autogenerated mock type for the GeospatialRepository type
type GeospatialRepository struct {
	mock.Mock
}

// FilteredDb provides a mock function with given fields: _a0
func (_m *GeospatialRepository) FilteredDb(_a0 model.GeospatialFilter) *gorm.DB {
	ret := _m.Called(_a0)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(model.GeospatialFilter) *gorm.DB); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *GeospatialRepository) Get(_a0 context.Context, _a1 model.GeospatialFilter) ([]model.Geospatial, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []model.Geospatial
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.GeospatialFilter) ([]model.Geospatial, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.GeospatialFilter) []model.Geospatial); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Geospatial)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.GeospatialFilter) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLevels provides a mock function with given fields: _a0
func (_m *GeospatialRepository) GetLevels(_a0 context.Context) ([]uint, error) {
	ret := _m.Called(_a0)

	var r0 []uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]uint, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []uint); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPaginate provides a mock function with given fields: _a0, _a1, _a2
func (_m *GeospatialRepository) GetPaginate(_a0 context.Context, _a1 model.GeospatialFilter, _a2 pagination.Param) ([]model.Geospatial, *pagination.Param, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []model.Geospatial
	var r1 *pagination.Param
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, model.GeospatialFilter, pagination.Param) ([]model.Geospatial, *pagination.Param, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.GeospatialFilter, pagination.Param) []model.Geospatial); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Geospatial)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.GeospatialFilter, pagination.Param) *pagination.Param); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*pagination.Param)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, model.GeospatialFilter, pagination.Param) error); ok {
		r2 = rf(_a0, _a1, _a2)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTypes provides a mock function with given fields: _a0
func (_m *GeospatialRepository) GetTypes(_a0 context.Context) ([]string, error) {
	ret := _m.Called(_a0)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertBulk provides a mock function with given fields: _a0, _a1
func (_m *GeospatialRepository) UpsertBulk(_a0 context.Context, _a1 []model.Geospatial) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []model.Geospatial) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewGeospatialRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewGeospatialRepository creates a new instance of GeospatialRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGeospatialRepository(t mockConstructorTestingTNewGeospatialRepository) *GeospatialRepository {
	mock := &GeospatialRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
