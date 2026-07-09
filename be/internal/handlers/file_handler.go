package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
)

// FileUpload handles file upload.
// @Summary Upload file
// @Description Upload a file to temporary storage
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Success 200 {object} dtos.UploadFileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/upload [post]
func (h *Handlers) FileUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		helpers.BadRequest(c, "File is required")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		helpers.BadRequest(c, "File extension is required")
		return
	}

	fileUUID := uuid.New().String()
	uploadDir := "storage/tmp"

	filePath, err := helpers.SaveUploadedFileWithOpts(c, "file", uploadDir, &helpers.SaveFileOptions{
		CustomName: fileUUID,
	})
	if err != nil {
		helpers.InternalServerError(c, fmt.Sprintf("Failed to upload file: %v", err))
		return
	}

	fileURL := helpers.GetFileURL(filePath)

	helpers.OK(c, "File uploaded successfully", dtos.UploadFileResponse{
		UUID: fileUUID,
		URL:  fileURL,
	})
}
