package handler

import (
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
}

func New(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Init(handlerConfig config.HandlerConfig) {
	mode := gin.ReleaseMode
	if handlerConfig.Server.Debug {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	h.router = gin.New()
	h.router.MaxMultipartMemory = int64(handlerConfig.Form.MaxSize.Bytes())

	h.router.Use(ginzap.Ginzap(h.logger, time.RFC3339, true))

	h.router.Use(cors.New(cors.Config{
		AllowOrigins:     handlerConfig.Server.Cors.AllowOrigins,
		AllowMethods:     handlerConfig.Server.Cors.AllowMethods,
		AllowHeaders:     handlerConfig.Server.Cors.AllowHeaders,
		AllowCredentials: handlerConfig.Server.Cors.AllowCredentials,
		ExposeHeaders:    handlerConfig.Server.Cors.ExposeHeaders,
		MaxAge:           handlerConfig.Server.Cors.MaxAge,
	}))

	hRest := rest.New(h.router, h.service, handlerConfig)
	hRest.Init()

	hSSE := sse.New(h.router, h.service)
	hSSE.Init()
}

func (h *Handler) Run(address string) error {
	return h.router.Run(address)
}
