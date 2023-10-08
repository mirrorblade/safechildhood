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
	router  *gin.Engine
	service *service.Service
	logger  *zap.Logger

	handlerConfig config.HandlerConfig
}

func New(service *service.Service, logger *zap.Logger, handlerConfig config.HandlerConfig) *Handler {
	return &Handler{
		service:       service,
		logger:        logger,
		handlerConfig: handlerConfig,
	}
}

func (h *Handler) Init() {
	mode := gin.ReleaseMode
	if h.handlerConfig.Server.Debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	h.router = gin.New()
	h.router.MaxMultipartMemory = int64(h.handlerConfig.Form.MaxSize.Bytes())

	h.router.Use(ginzap.Ginzap(h.logger, time.RFC3339, true))

	h.router.Use(cors.New(cors.Config{
		AllowOrigins:     h.handlerConfig.Server.Cors.AllowOrigins,
		AllowMethods:     h.handlerConfig.Server.Cors.AllowMethods,
		AllowHeaders:     h.handlerConfig.Server.Cors.AllowHeaders,
		AllowCredentials: h.handlerConfig.Server.Cors.AllowCredentials,
		ExposeHeaders:    h.handlerConfig.Server.Cors.ExposeHeaders,
		MaxAge:           h.handlerConfig.Server.Cors.MaxAge,
	}))

	h.initConfig()

	hRest := rest.New(h.router, h.service, h.handlerConfig.Form.MaxPhotos)
	hRest.Init()

	hSSE := sse.New(h.router, h.service)
	hSSE.Init()
}

func (h *Handler) Run(address string) error {
	return h.router.Run(address)
}

func (h *Handler) initConfig() {
	h.router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"map":        h.handlerConfig.Map,
			"complaints": h.handlerConfig.Complaints,
		})
	})
}
