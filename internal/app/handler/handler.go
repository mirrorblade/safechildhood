package handler

import (
	"safechildhood/internal/app/handler/rest"
	"safechildhood/internal/app/handler/sse"
	"safechildhood/internal/app/service"
	"time"

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

func (h *Handler) Init() {
	h.router = gin.Default()
	h.router.MaxMultipartMemory = 15 << 20

	h.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	hRest := rest.New(h.router, h.service)
	hRest.Init()

	hSSE := sse.New(h.router, h.service)
	hSSE.Init()
}

func (h *Handler) Run(address string) error {
	return h.router.Run(address)
}
