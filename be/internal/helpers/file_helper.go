package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var defaultAllowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".jfif": true,
}

const defaultMaxSizeMB = 5

type SaveFileOptions struct {
	AllowedExts []string
	MaxSizeMB   int
	CustomName  string
}

type FileMetadata struct {
	UUID      string
	Extension string
	SizeBytes int64
	SizeMB    float64
	FullPath  string
}

func SaveUploadedFileWithOpts(c *gin.Context, fieldName string, uploadDir string, opts *SaveFileOptions) (string, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	allowedExts := defaultAllowedExts
	maxSize := int64(defaultMaxSizeMB) * 1024 * 1024

	if opts != nil {
		if len(opts.AllowedExts) > 0 {
			allowedExts = make(map[string]bool)
			for _, ext := range opts.AllowedExts {
				allowedExts[strings.ToLower(ext)] = true
			}
		}
		if opts.MaxSizeMB > 0 {
			maxSize = int64(opts.MaxSizeMB) * 1024 * 1024
		}
	}

	if header.Size > maxSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %dMB", maxSize/1024/1024)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("file type %s is not allowed", ext)
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	var fileName string
	if opts != nil && opts.CustomName != "" {
		fileName = opts.CustomName + ext
	} else {
		randomStr, err := GenerateRandomString(4)
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		fileName = fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomStr, ext)
	}

	filePath := filepath.Join(uploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	normalizedDir := strings.ReplaceAll(filepath.ToSlash(uploadDir), "\\", "/")
	return fmt.Sprintf("%s/%s", normalizedDir, fileName), nil
}

func SaveUploadedFile(c *gin.Context, fieldName string, uploadDir string) (string, error) {
	return SaveUploadedFileWithOpts(c, fieldName, uploadDir, nil)
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func GetFileURL(path string) string {
	if path == "" {
		return ""
	}
	baseURL := GetEnv("APP_URL", "http://localhost:8080")
	normalizedPath := strings.ReplaceAll(path, "\\", "/")
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), normalizedPath)
}

func CopyFile(fileUUID, srcDir, destDir string) (string, error) {
	if fileUUID == "" {
		return "", fmt.Errorf("file UUID is empty")
	}

	files, err := filepath.Glob(filepath.Join(srcDir, fileUUID+".*"))
	if err != nil {
		return "", fmt.Errorf("failed to search file: %w", err)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("file not found for uuid: %s", fileUUID)
	}

	srcPath := files[0]
	ext := filepath.Ext(srcPath)
	destPath := filepath.Join(destDir, fileUUID+ext)

	if err := copyFileInternal(srcPath, destPath); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	normalizedPath := strings.ReplaceAll(filepath.ToSlash(destPath), "\\", "/")
	return normalizedPath, nil
}

func copyFileInternal(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func MoveFile(fileUUID, srcDir, destDir string) (string, error) {
	if fileUUID == "" {
		return "", fmt.Errorf("file UUID is empty")
	}

	files, err := filepath.Glob(filepath.Join(srcDir, fileUUID+".*"))
	if err != nil {
		return "", fmt.Errorf("failed to search file: %w", err)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("file not found for uuid: %s", fileUUID)
	}

	srcPath := files[0]
	ext := filepath.Ext(srcPath)
	destPath := filepath.Join(destDir, fileUUID+ext)

	if err := copyFileInternal(srcPath, destPath); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	if err := DeleteFile(srcPath); err != nil {
		return "", fmt.Errorf("failed to delete source file: %w", err)
	}

	normalizedPath := strings.ReplaceAll(filepath.ToSlash(destPath), "\\", "/")
	return normalizedPath, nil
}

func GetFileMetadata(fileUUID, dir string) (*FileMetadata, error) {
	if fileUUID == "" {
		return nil, fmt.Errorf("file UUID is empty")
	}

	files, err := filepath.Glob(filepath.Join(dir, fileUUID+".*"))
	if err != nil {
		return nil, fmt.Errorf("failed to search file: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found for uuid: %s", fileUUID)
	}

	srcPath := files[0]
	info, err := os.Stat(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(srcPath))
	sizeMB := float64(info.Size()) / 1024.0 / 1024.0

	return &FileMetadata{
		UUID:      fileUUID,
		Extension: ext,
		SizeBytes: info.Size(),
		SizeMB:    sizeMB,
		FullPath:  srcPath,
	}, nil
}
