package helpers

import (
	"regexp"
	"strings"
)

// Deklarasi regex di luar fungsi agar dikompilasi 1x saja saat aplikasi berjalan
var (
	nonAlphaNumRegex = regexp.MustCompile(`[^a-z0-9]+`)
	trimDashRegex    = regexp.MustCompile(`^-+|-+$`)
)

func Slugify(text string) string {
	slug := strings.ToLower(text)

	// Ganti semua urutan karakter non-alfanumerik dengan satu strip
	slug = nonAlphaNumRegex.ReplaceAllString(slug, "-")

	// Hapus strip di awal dan di akhir string (jika ada)
	return trimDashRegex.ReplaceAllString(slug, "")
}
