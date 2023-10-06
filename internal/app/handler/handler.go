package handler

import (
	"safechildhood/internal/app/config"
	"safechildhood/internal/app/handler/rest"
	"safechildhood/internal/app/handler/sse"
	"safechildhood/internal/app/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	router  *gin.Engine
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init(handlerConfig config.HandlerConfig) {
	h.router = gin.Default()
	h.router.MaxMultipartMemory = int64(handlerConfig.Form.MaxSize.Bytes())

	h.router.Use(cors.New(cors.Config{
		AllowOrigins:     handlerConfig.Cors.AllowOrigins,
		AllowMethods:     handlerConfig.Cors.AllowMethods,
		AllowHeaders:     handlerConfig.Cors.AllowHeaders,
		AllowCredentials: handlerConfig.Cors.AllowCredentials,
		MaxAge:           handlerConfig.Cors.MaxAge,
	}))

	hRest := rest.New(h.router, h.service, handlerConfig)
	hRest.Init()

	hSSE := sse.New(h.router, h.service)
	hSSE.Init()
}

func (h *Handler) Run(address string) error {
	return h.router.Run(address)
}
