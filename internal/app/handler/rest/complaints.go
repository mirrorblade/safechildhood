package rest

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"safechildhood/internal/app/domain"
	"safechildhood/pkg/converter"
	"safechildhood/pkg/storage"
	"slices"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid/v5"
)

var commonImageMimeTypes = []string{"image/avif", "image/bmp", "image/gif", "image/jpeg", "image/png", "image/tiff", "image/webp"}

func (h *Handler) initComplaints() {
	basicAuth := gin.BasicAuth(gin.Accounts{
		"admin": "123456789",
	})

	complaint := h.router.Group("/complaints")
	{
		complaint.GET("/:id", basicAuth, h.getComplaint)
		complaint.POST("/", h.createComplaint)
		complaint.DELETE("/:id", basicAuth, h.deleteComplaint)
	}
}

type createComplaintBody struct {
	Coordinates      string                  `form:"coordinates" json:"coordinates" binding:"required,max=32"`
	ShortDescription string                  `form:"short_description" json:"short_description" binding:"required,max=100"`
	Description      string                  `form:"description" json:"description" binding:"required,max=5000"`
	Photos           []*multipart.FileHeader `form:"photos" json:"photos"`
}

func (h *Handler) getComplaint(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 36 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"body":    nil,
		})

		return
	}

	if _, err := uuid.FromString(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"body":    nil,
		})

		return
	}

	complaint, err := h.service.Complaints.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "bad request",
				"body":    nil,
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "server error",
			"body":    nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": nil,
		"body":    complaint,
	})
}

func (h *Handler) createComplaint(c *gin.Context) {
	var bodyComplaint createComplaintBody

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": "bad asdsadsafasf",
		"body":    nil,
	})

	return

	if err := c.ShouldBindWith(&bodyComplaint, binding.FormMultipart); err != nil {
		if bodyComplaint.Coordinates == "" || bodyComplaint.ShortDescription == "" || bodyComplaint.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "bad request",
				"body":    nil,
			})

			return
		}
	}

	if _, err := converter.StringToCoordinates(bodyComplaint.Coordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"body":    nil,
		})

		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "server error",
			"body":    nil,
		})

		return
	}

	complaint := new(domain.Complaint)
	complaint.ID = id
	complaint.Coordinates = bodyComplaint.Coordinates
	complaint.ShortDescription = h.textSanitazer.Sanitize(bodyComplaint.ShortDescription)
	complaint.Description = h.textSanitazer.Sanitize(bodyComplaint.Description)
	complaint.CreatedAt = time.Now()

	if len(bodyComplaint.Photos) != 0 {
		if len(bodyComplaint.Photos) > h.handlerConfig.Form.MaxPhotos {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": domain.ErrTooManyPhotos.Error() + fmt.Sprintf(" (max photos count: %d)", h.handlerConfig.Form.MaxPhotos),
				"body":    nil,
			})

			return
		}

		folderId := h.service.GetSavedFolderId(bodyComplaint.Coordinates)
		if folderId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "bad request",
				"body":    nil,
			})

			return
		}

		complaintDir, err := h.service.Storage.Create(c.Request.Context(), storage.GoogleDriveParameters{
			Name:       id.String(),
			ObjectMode: storage.FOLDER,
			ParentId:   folderId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "server error",
				"body":    nil,
			})

			return
		}

		for _, photo := range bodyComplaint.Photos {
			file, err := photo.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "server error",
					"body":    nil,
				})

				return
			}

			defer file.Close()

			mimeType, err := mimetype.DetectReader(file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "bad request",
					"body":    nil,
				})

				return
			}

			if !slices.Contains(commonImageMimeTypes, mimeType.Extension()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "bad request",
					"body":    nil,
				})

				return
			}

			if _, err := h.service.Storage.Create(c.Request.Context(), storage.GoogleDriveParameters{
				Name:                   photo.Filename,
				ObjectMode:             storage.FILE,
				Content:                file,
				ParentId:               complaintDir.Id,
				SkipAlreadyExistsCheck: true,
			}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "server error",
					"body":    nil,
				})

				return
			}
		}

		complaint.PhotosPath = fmt.Sprintf("/%s/%s", bodyComplaint.Coordinates, id)
	}

	if err := h.service.Complaints.Create(c.Request.Context(), *complaint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "server error",
			"body":    nil,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "complaint successfully created",
		"body":    nil,
	})
}

func (h *Handler) deleteComplaint(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 36 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"body":    nil,
		})

		return
	}

	if _, err := uuid.FromString(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"body":    nil,
		})

		return
	}

	if err := h.service.Complaints.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
				"body":    nil,
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "server error",
			"body":    nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user was successfully deleted",
		"body":    nil,
	})
}
