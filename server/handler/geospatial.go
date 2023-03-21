package handler

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-geospatial/domain/model"
	"github.com/si-bas/go-rest-geospatial/pkg/logger"
	"github.com/si-bas/go-rest-geospatial/pkg/logger/tag"
	"github.com/si-bas/go-rest-geospatial/shared/helper/pagination"
	"github.com/si-bas/go-rest-geospatial/shared/helper/response"
	"github.com/twpayne/go-geom/encoding/geojson"
)

func validateGeospatialFilter(query model.GeospatialFilterParams) (*model.GeospatialFilter, error) {
	filter := model.GeospatialFilter{
		Name: query.Name,
	}

	if query.LatLng != "" {
		errMsg := "latlng must contain two float values, divided by commas"

		latLng := strings.Split(query.LatLng, ",")
		if len(latLng) != 2 {
			return nil, errors.New(errMsg)
		}

		fLat, err := strconv.ParseFloat(latLng[0], 64)
		if err != nil {
			return nil, errors.New(errMsg)
		}

		fLng, err := strconv.ParseFloat(latLng[1], 64)
		if err != nil {
			return nil, errors.New(errMsg)
		}

		filter.Lat = fLat
		filter.Lng = fLng

		filter.Nested = query.Nested
	}

	if query.Types != "" {
		filter.Types = strings.Split(query.Types, ",")
	}

	if query.Levels != "" {
		errMsg := "levels must contain one or more integer values, divided by commas"
		levels := strings.Split(query.Levels, ",")

		for _, str := range levels {
			uintVal, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return nil, errors.New(errMsg)
			}
			filter.Levels = append(filter.Levels, uint(uintVal))
		}
	}

	if query.ExcludedIds != "" {
		errMsg := "excluded ids must contain one or more integer values, divided by commas"
		excludedIds := strings.Split(query.ExcludedIds, ",")

		for _, str := range excludedIds {
			uintVal, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return nil, errors.New(errMsg)
			}
			filter.ExcludedIds = append(filter.ExcludedIds, uint(uintVal))
		}
	}

	if query.ParentIds != "" {
		errMsg := "parent ids must contain one or more integer values, divided by commas"
		parentIds := strings.Split(query.ParentIds, ",")

		for _, str := range parentIds {
			uintVal, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return nil, errors.New(errMsg)
			}
			filter.ParentIds = append(filter.ParentIds, uint(uintVal))
		}
	}

	return &filter, nil
}

func (h *Handler) GeospatialList(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse(ctx)

	var query model.GeospatialFilterParams
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Warn(ctx, "failed to bindQuery", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	filter, err := validateGeospatialFilter(query)
	if err != nil {
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	if filter.Lat != 0 && filter.Lng != 0 {
		data, err := h.geospatialService.List(ctx, *filter)
		if err != nil {
			c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
			return
		}

		if filter.Nested {
			var rootLevel uint = ^uint(0)
			for _, g := range data {
				if g.Level < rootLevel {
					rootLevel = g.Level
				}
			}
			root := h.geospatialService.BuildTree(data, rootLevel)
			c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(root))
			return
		}

		c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(data))
		return
	}

	var sortBys []pagination.ParamSort
	if len(query.Sort) > 0 {
		for k, v := range query.Sort {
			sortBys = append(sortBys, pagination.ParamSort{
				Column: k,
				Order:  v,
			})
		}
	}

	data, meta, err := h.geospatialService.ListPaginate(ctx, *filter, pagination.Param{
		Limit: query.Limit,
		Page:  query.Page,
		Sort:  sortBys,
	})
	if err != nil {
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(data).SetMeta(meta))
}

func (h *Handler) GeospatialTypes(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse(ctx)

	data, err := h.geospatialService.GetTypes(ctx)
	if err != nil {
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}
	if data == nil {
		data = make([]string, 0)
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(data))
}

func (h *Handler) GeospatialLevels(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse(ctx)

	data, err := h.geospatialService.GetLevels(ctx)
	if err != nil {
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}
	if data == nil {
		data = make([]uint, 0)
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(data))
}

func (h *Handler) GeospatialImport(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse(ctx)

	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn(ctx, "failed to get file from request", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	geojsonFile, err := file.Open()
	if err != nil {
		logger.Warn(ctx, "failed to open file from request", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}
	defer geojsonFile.Close()

	geojsonData, err := io.ReadAll(geojsonFile)
	if err != nil {
		logger.Warn(ctx, "failed to read file from request", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	var features geojson.FeatureCollection
	if err := json.Unmarshal(geojsonData, &features); err != nil {
		logger.Warn(ctx, "failed to Unmarshal", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	if err := h.geospatialService.CreateFromFeatureCollection(ctx, &features); err != nil {
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result)
}
