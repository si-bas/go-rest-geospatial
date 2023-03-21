package handler

import "github.com/si-bas/go-rest-geospatial/service"

type Handler struct {
	geospatialService service.GeospatialService
}

func New(
	geospatialService service.GeospatialService,
) *Handler {
	return &Handler{
		geospatialService: geospatialService,
	}
}
