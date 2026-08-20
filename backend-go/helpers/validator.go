package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TranslateErrorMessage menangani error validasi dan database menjadi map yang lebih ramah
func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Handle validasi dari validator.v10
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s wajib diisi", field)
			case "email":
				errorsMap[field] = "Format email tidak valid"
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s sudah terdaftar", field)
			case "min":
				errorsMap[field] = fmt.Sprintf("%s minimal %s karakter", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s maksimal %s karakter", field, fieldError.Param())
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s harus berupa angka", field)
			default:
				errorsMap[field] = "Nilai tidak valid"
			}
		}
	}

	// Handle GORM error: Duplicate entry
	if err != nil {
		if strings.Contains(err.Error(), "sudah terdaftar") {
			field := extractDuplicateField(err.Error())
			if field != "" {
				errorsMap[field] = fmt.Sprintf("%s sudah terdaftar", field)
			} else {
				errorsMap["Error"] = "Sudah terdaftar"
			}
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["Error"] = "Data tidak ditemukan"
		}
	}

	return errorsMap
}

// IsDuplicateEntryError mengecek apakah error adalah duplicate entry
func IsDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Sudah terdaftar")
}

// extractDuplicateField mencoba mengekstrak nama kolom dari error "Duplicate entry"
func extractDuplicateField(errMsg string) string {
	// Contoh error MySQL: Error 1062: Duplicate entry 'test@example.com' for key 'users.email'
	// Kita ambil bagian setelah 'for key' lalu extract nama field
	re := regexp.MustCompile(`for key '(\w+\.)?(\w+)'`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) == 3 {
		// Hasilkan kapitalisasi nama field
		return capitalize(matches[2])
	}
	return ""
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
