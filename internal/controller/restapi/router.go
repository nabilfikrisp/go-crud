// Package restapi configures HTTP routing and REST API handlers.
package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nabilfikrisp/go-crud/config"
	_ "github.com/nabilfikrisp/go-crud/docs" // Swagger Docs
	"github.com/nabilfikrisp/go-crud/internal/controller/restapi/middleware"
	v1 "github.com/nabilfikrisp/go-crud/internal/controller/restapi/v1"
	"github.com/nabilfikrisp/go-crud/internal/usecase"
	"github.com/nabilfikrisp/go-crud/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
//
//	@title       Go Clean Crud API
//	@description Crud clean architecture template with contact, and user management
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /v1
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
func NewRouter(engine *gin.Engine, cfg *config.Config, c usecase.Contact, l logger.Interface) {
	// Options
	engine.Use(middleware.Logger(l))
	engine.Use(middleware.Recovery(l))

	// Swagger
	if cfg.Swagger.Enabled {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Routers
	apiV1Group := engine.Group("/v1")
	{
		v1.NewRoutes(*apiV1Group, c, l)
	}
}
