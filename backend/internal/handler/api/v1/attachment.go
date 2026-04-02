package v1

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	uploadDir string
}

func NewAttachmentHandler() *AttachmentHandler {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	os.MkdirAll(uploadDir, 0755)
	return &AttachmentHandler{uploadDir: uploadDir}
}

type UploadResponse struct {
	ID        uint      `json:"id"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// 验证文件大小 (10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large, max 10MB"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(h.uploadDir, newFilename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		ID:        1,
		FileName:  file.Filename,
		FileSize:  file.Size,
		MimeType:  file.Header.Get("Content-Type"),
		Path:      filePath,
		CreatedAt: time.Now(),
	})
}

func (h *AttachmentHandler) Download(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	c.JSON(http.StatusOK, gin.H{
		"id":          id,
		"file_name":   "test.pdf",
		"download_url": fmt.Sprintf("/api/v1/attachments/%d/download", id),
	})
}

func (h *AttachmentHandler) Preview(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	c.JSON(http.StatusOK, gin.H{
		"id":         id,
		"preview_url": fmt.Sprintf("/api/v1/attachments/%d/preview", id),
	})
}
