package helpers

import (
	"backend-go/structs"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadConfig struct {
	File           *multipart.FileHeader // File yang akan diupload
	AllowedTypes   []string              // Ekstensi file yang diperbolehkan (bisa ".jpg" atau "jpg")
	MaxSize        int64                 // Ukuran maksimum file (dalam byte)
	DestinationDir string                // Folder tujuan penyimpanan
}

type UploadResult struct {
	FileName string
	FilePath string
	Error    error
	Response *structs.ErrorResponse
}

// Untuk nama file yang sudah di-slugify (opsional jika ingin nama asli di-slugify alih-alih UUID)
func SlugifyFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	slugBase := Slugify(base)
	return slugBase + ext
}

func UploadFile(c *gin.Context, config UploadConfig) UploadResult {
	// 1. Cek apakah file ada
	if config.File == nil {
		err := errors.New("file is required")
		return UploadResult{
			Error: err,
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File is required",
				Errors:  map[string]string{"file": "No file was uploaded"},
			},
		}
	}

	// 2. Cek ukuran file (Gunakan float64 agar akurat jika < 1MB)
	if config.File.Size > config.MaxSize {
		maxMB := float64(config.MaxSize) / (1024 * 1024)
		err := fmt.Errorf("file size exceeds maximum limit of %.2fMB", maxMB)
		return UploadResult{
			Error: err,
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File too large",
				Errors:  map[string]string{"file": fmt.Sprintf("Maximum file size is: %.2fMB", maxMB)},
			},
		}
	}

	// 3. Normalisasi ekstensi file dan AllowedTypes agar aman dari titik/huruf besar-kecil
	ext := strings.ToLower(filepath.Ext(config.File.Filename))
	cleanExt := strings.TrimPrefix(ext, ".")

	var cleanAllowedTypes []string
	for _, t := range config.AllowedTypes {
		cleanAllowedTypes = append(cleanAllowedTypes, strings.TrimPrefix(strings.ToLower(t), "."))
	}

	if !slices.Contains(cleanAllowedTypes, cleanExt) {
		err := fmt.Errorf("file type %s is not allowed", ext)
		return UploadResult{
			Error: err,
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Invalid file type",
				Errors:  map[string]string{"file": fmt.Sprintf("Allowed file types: %v", config.AllowedTypes)},
			},
		}
	}

	// 4. Generate UUID sebagai nama file
	uuidName := uuid.New().String()
	filename := uuidName + ext

	// filepath.ToSlash mengubah path "\\" Windows menjadi "/" agar aman untuk URL/database
	filePath := filepath.ToSlash(filepath.Join(config.DestinationDir, filename))

	// 5. Buat folder tujuan jika belum ada
	if err := os.MkdirAll(config.DestinationDir, 0755); err != nil {
		return UploadResult{
			Error: err,
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to create upload directory",
				Errors:  map[string]string{"system": err.Error()},
			},
		}
	}

	// 6. Simpan file ke folder tujuan
	if err := c.SaveUploadedFile(config.File, filePath); err != nil {
		return UploadResult{
			Error: err,
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to save file",
				Errors:  map[string]string{"file": err.Error()},
			},
		}
	}

	// Kembalikan hasil upload sukses
	return UploadResult{
		FileName: filename,
		FilePath: filePath,
	}
}
