package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/si-bas/go-rest-geospatial/domain/model"
	"github.com/si-bas/go-rest-geospatial/shared/helper/pagination"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
	"gorm.io/gorm"
)

type GeospatialRepository interface {
	FilteredDb(model.GeospatialFilter) *gorm.DB

	Get(context.Context, model.GeospatialFilter) ([]model.Geospatial, error)
	GetPaginate(context.Context, model.GeospatialFilter, pagination.Param) ([]model.Geospatial, *pagination.Param, error)
	GetTypes(context.Context) ([]string, error)
	GetLevels(context.Context) ([]uint, error)
	UpsertBulk(context.Context, []model.Geospatial) error
}

type geospatialImpl struct {
	db *gorm.DB
}

func NewGeospatialRepository(db *gorm.DB) GeospatialRepository {
	return &geospatialImpl{
		db: db,
	}
}

func (r *geospatialImpl) FilteredDb(filter model.GeospatialFilter) *gorm.DB {
	chain := r.db.Model(&model.Geospatial{})

	if filter.Name != "" {
		chain.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Levels != nil && len(filter.Levels) > 0 {
		chain.Where("level IN (?)", filter.Levels)
	}

	if filter.Types != nil && len(filter.Types) > 0 {
		chain.Where("type IN (?)", filter.Types)
	}

	if filter.ExcludedIds != nil && len(filter.ExcludedIds) > 0 {
		chain.Where("id NOT IN (?)", filter.ExcludedIds)
	}

	if filter.ParentIds != nil && len(filter.ParentIds) > 0 {
		chain.Where("parent_gadm_id IN (SELECT gadm_id FROM geospatial WHERE id IN (?))", filter.ParentIds)
	}

	if filter.Lat != 0 && filter.Lng != 0 {
		point := geom.NewPointFlat(geom.XY, []float64{filter.Lng, filter.Lat})
		wktString, err := wkt.Marshal(point)
		if err == nil {
			chain.Where("ST_Contains(geometry, ST_GeomFromText(?))", wktString)
		}
	}

	return chain
}

func (r *geospatialImpl) Get(ctx context.Context, filter model.GeospatialFilter) ([]model.Geospatial, error) {
	var geospatials []model.Geospatial

	if err := r.FilteredDb(filter).Order("level ASC").Find(&geospatials).Error; err != nil {
		return nil, err
	}

	return geospatials, nil
}

func (r *geospatialImpl) GetPaginate(ctx context.Context, filter model.GeospatialFilter, param pagination.Param) ([]model.Geospatial, *pagination.Param, error) {
	var geospatials []model.Geospatial

	filteredDb := r.FilteredDb(filter)

	if err := filteredDb.Scopes(pagination.Paginate(model.Geospatial{}, &param, filteredDb)).Find(&geospatials).Error; err != nil {
		return nil, nil, err
	}

	return geospatials, &param, nil
}

func (r *geospatialImpl) GetTypes(ctx context.Context) ([]string, error) {
	var geospatials []model.Geospatial
	if err := r.db.Model(&model.Geospatial{}).Select("type").Group("type").Find(&geospatials).Error; err != nil {
		return nil, err
	}

	var types []string
	for _, g := range geospatials {
		types = append(types, g.Type)
	}

	return types, nil
}

func (r *geospatialImpl) GetLevels(ctx context.Context) ([]uint, error) {
	var geospatials []model.Geospatial
	if err := r.db.Model(&model.Geospatial{}).Select("level").Group("level").Find(&geospatials).Error; err != nil {
		return nil, err
	}

	var levels []uint
	for _, g := range geospatials {
		levels = append(levels, g.Level)
	}

	return levels, nil
}

func (r *geospatialImpl) UpsertBulk(ctx context.Context, geospatials []model.Geospatial) error {
	var values []interface{}
	var placeholders []string
	for _, g := range geospatials {
		values = append(values, g.GadmID, g.ParentGadmID, g.Name, g.Type, g.Level, g.Geometry)
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ST_GeomFromText(?))")
	}

	query := fmt.Sprintf("INSERT INTO geospatial (gadm_id, parent_gadm_id, name, type, level, geometry) VALUES %s ON DUPLICATE KEY UPDATE gadm_id=VALUES(gadm_id), parent_gadm_id=VALUES(parent_gadm_id), name=VALUES(name), type=VALUES(type), level=VALUES(level), geometry=VALUES(geometry)", strings.Join(placeholders, ", "))
	result := r.db.Exec(query, values...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
