package handler

import (
	"net/http"
	"safechildhood/internal/app/config"
	"safechildhood/internal/app/handler/rest"
	"safechildhood/internal/app/handler/sse"
	"safechildhood/internal/app/service"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	router *gin.Engine

	service *service.Service
	logger  *zap.Logger
}

func New(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Init(config *config.Config) {
	mode := gin.ReleaseMode
	if config.Server.Debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	h.router = gin.New()
	h.router.MaxMultipartMemory = int64(config.Form.MaxSize.Bytes())

	h.router.Use(ginzap.Ginzap(h.logger, time.RFC3339, true))

	h.router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Server.Cors.AllowOrigins,
		AllowMethods:     config.Server.Cors.AllowMethods,
		AllowHeaders:     config.Server.Cors.AllowHeaders,
		AllowCredentials: config.Server.Cors.AllowCredentials,
		ExposeHeaders:    config.Server.Cors.ExposeHeaders,
		MaxAge:           config.Server.Cors.MaxAge,
	}))

	h.initConfig(config)
	h.initRest(config)
	h.initSSE(config)

}

func (h *Handler) Run(address string) error {
	return h.router.Run(address)
}

func (h *Handler) initRest(config *config.Config) {
	hRest := rest.New(h.service, config.Form.MaxPhotos)
	gRest := h.router.Group("/api")
	{
		hRest.Init(gRest)
	}
}

func (h *Handler) initSSE(config *config.Config) {
	hSSE := sse.New(h.service)
	gSSE := h.router.Group("/sse", sse.HeadersMiddleware())
	{
		hSSE.Init(gSSE)
	}
}

func (h *Handler) initConfig(config *config.Config) {
	h.router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"map": config.Map,
		})
	})
}
