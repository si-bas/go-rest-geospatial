package service

import (
	"context"
	"fmt"

	"github.com/si-bas/go-rest-geospatial/domain/model"
	"github.com/si-bas/go-rest-geospatial/domain/repository"
	"github.com/si-bas/go-rest-geospatial/pkg/logger"
	"github.com/si-bas/go-rest-geospatial/shared"
	"github.com/si-bas/go-rest-geospatial/shared/helper/pagination"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
	"gorm.io/gorm"
)

type GeospatialService interface {
	List(context.Context, model.GeospatialFilter) ([]model.Geospatial, error)
	ListPaginate(context.Context, model.GeospatialFilter, pagination.Param) ([]model.Geospatial, *pagination.Param, error)
	GetTypes(context.Context) ([]string, error)
	GetLevels(context.Context) ([]uint, error)
	CreateFromFeatureCollection(context.Context, *geojson.FeatureCollection) error
	BuildTree([]model.Geospatial, uint) *model.Geospatial
}

type geospatialImpl struct {
	geospatialRepo repository.GeospatialRepository
}

func NewGeospatialService(geospatialRepo repository.GeospatialRepository) GeospatialService {
	return &geospatialImpl{
		geospatialRepo: geospatialRepo,
	}
}

func (s *geospatialImpl) List(ctx context.Context, filter model.GeospatialFilter) ([]model.Geospatial, error) {
	geospatials, err := s.geospatialRepo.Get(ctx, filter)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error(ctx, "failed to get geospatial data with filter", err)
		}

		return nil, err
	}

	return geospatials, nil
}

func (s *geospatialImpl) ListPaginate(ctx context.Context, filter model.GeospatialFilter, query pagination.Param) ([]model.Geospatial, *pagination.Param, error) {
	geospatials, meta, err := s.geospatialRepo.GetPaginate(ctx, filter, query)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error(ctx, "failed to get geospatial data with filter", err)
		}

		return nil, nil, err
	}

	return geospatials, meta, nil
}

func (s *geospatialImpl) GetTypes(ctx context.Context) ([]string, error) {
	types, err := s.geospatialRepo.GetTypes(ctx)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error(ctx, "failed to get geospatial types", err)
		}

		return nil, err
	}

	return types, nil
}

func (s *geospatialImpl) GetLevels(ctx context.Context) ([]uint, error) {
	levels, err := s.geospatialRepo.GetLevels(ctx)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error(ctx, "failed to get geospatial types", err)
		}

		return nil, err
	}

	return levels, nil
}

func (s *geospatialImpl) CreateFromFeatureCollection(ctx context.Context, fc *geojson.FeatureCollection) error {
	var geospatials []model.Geospatial
	for _, f := range fc.Features {

		// Get level
		var level int
		for i := 0; i <= 4; i++ {
			if f.Properties[fmt.Sprintf("GID_%d", i)] == nil {
				break
			}

			level = i
		}

		wktEncoder := wkt.NewEncoder()
		mpStr, _ := wktEncoder.Encode(f.Geometry)

		geospatial := model.Geospatial{
			GadmID:   f.Properties[fmt.Sprintf("GID_%d", level)].(string),
			Name:     f.Properties["COUNTRY"].(string),
			Type:     "Country",
			Level:    uint(level + 1),
			Geometry: mpStr,
		}

		if level > 0 {
			geospatial.ParentGadmID = f.Properties[fmt.Sprintf("GID_%d", level-1)].(string)
			geospatial.Name = f.Properties[fmt.Sprintf("NAME_%d", level)].(string)
			geospatial.Type = f.Properties[fmt.Sprintf("TYPE_%d", level)].(string)
		}

		geospatials = append(geospatials, geospatial)
	}

	chunkSize := 1000
	geospatialChunks := shared.ChunkGeospatialData(geospatials, chunkSize)

	var result []error
	ch := make(chan error, len(geospatialChunks))
	for _, c := range geospatialChunks {
		go func(ctx context.Context, data []model.Geospatial, result chan<- error) {
			if err := s.geospatialRepo.UpsertBulk(ctx, data); err != nil {
				result <- err
				return
			}
			result <- nil
		}(ctx, c, ch)
	}
	for i := 0; i < len(geospatialChunks); i++ {
		res := <-ch
		if res != nil {
			result = append(result, res)
		}
	}
	close(ch)

	if len(result) > 0 {
		for _, err := range result {
			return err
		}
	}

	return nil
}

func (s *geospatialImpl) BuildTree(geos []model.Geospatial, rootLevel uint) *model.Geospatial {
	// Find the root node with the smallest level
	var root *model.Geospatial
	for i := range geos {
		if geos[i].Level == rootLevel {
			root = &geos[i]
			break
		}
	}
	if root == nil {
		return nil
	}
	// Recursively build the tree from the root node
	for i := range geos {
		if geos[i].Level == rootLevel+1 && geos[i].ParentGadmID == root.GadmID {
			child := s.BuildTree(geos, rootLevel+1)
			if child != nil {
				root.Child = child
			}
		}
	}
	return root
}
