package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-geospatial/config"
	"github.com/si-bas/go-rest-geospatial/domain/repository"
	"github.com/si-bas/go-rest-geospatial/pkg/gorm"
	"github.com/si-bas/go-rest-geospatial/pkg/logger"
	"github.com/si-bas/go-rest-geospatial/server/handler"
	"github.com/si-bas/go-rest-geospatial/server/middleware"
	"github.com/si-bas/go-rest-geospatial/service"
	"github.com/si-bas/go-rest-geospatial/shared/constant"
)

type HTTPServer struct {
}

// New to instantiate HTTPServer
func New() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start() {
	h := initHandler()

	if config.Config.App.Env == constant.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.InjectContext())
	router.GET("/healthcheck", h.HealthCheck)

	groupV1 := router.Group("/v1")

	groupV1.GET("/q", h.GeospatialList)
	groupV1.GET("/types", h.GeospatialTypes)
	groupV1.GET("/levels", h.GeospatialLevels)
	groupV1.POST("/import", h.GeospatialImport)

	err := router.Run(fmt.Sprintf(":%d", config.Config.App.Port))
	if err != nil {
		logger.Error(context.Background(), "failed to run router", err)
	}
}

func initHandler() *handler.Handler {
	var err error
	config.TimeLocation, err = time.LoadLocation(config.Config.App.Timezone)
	if err != nil {
		panic("error set timezone, err=" + err.Error())
	}

	logger.InitLogger()

	// TODO: init DB
	db := gorm.ConnectDB()

	// TODO: init repositories
	geospatialRepo := repository.NewGeospatialRepository(db)

	// TODO: init pkgs

	// TODO: init services
	geospatialService := service.NewGeospatialService(geospatialRepo)

	return handler.New(
		geospatialService,
	)
}
