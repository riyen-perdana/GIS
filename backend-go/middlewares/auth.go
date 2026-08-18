package middlewares

import (
	"backend-go/config"
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Prefix standar pada header Authorization: "Authorization: Bearer <token>"
const bearerPrefix = "Bearer "

// jwtKey diinisialisasi SEKALI saat package di-load, bukan tiap request.
//
// Catatan: karena ini variabel package-level, mengubah .env saat aplikasi
// sedang berjalan tidak berpengaruh sampai aplikasi di-restart.
var jwtKey = mustJWTKey()

// mustJWTKey sengaja menghentikan aplikasi kalau JWT_KEY belum di-set.
//
// Ini perubahan penting dari versi lama yang memakai fallback "secret_key":
// fallback seperti itu membuat aplikasi TETAP JALAN dengan kunci yang bisa
// ditebak siapa pun yang pernah membaca repo ini — dan tidak ada yang sadar
// sampai ada yang menyalahgunakannya.
//
// Lebih baik gagal keras saat startup (developer langsung tahu)
// daripada gagal diam-diam di production.
func mustJWTKey() []byte {
	key := config.GetEnv("JWT_KEY", "")
	if key == "" {
		log.Fatal("[FATAL] JWT_KEY belum di-set. Tambahkan JWT_KEY di file .env sebelum menjalankan aplikasi.")
	}
	return []byte(key)
}

// Auth memastikan request membawa JWT yang sah, lalu menaruh identitas
// user (username) ke dalam context supaya bisa dipakai middleware dan
// controller berikutnya.
//
// Middleware ini SENGAJA tidak menyentuh database. Tugasnya hanya
// menjawab "kamu siapa?", bukan "kamu boleh apa?".
// Urusan hak akses ada di middleware Permission.
//
// Contoh pemakaian:
//
//	api := router.Group("/api", middlewares.Auth())
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := extractBearerToken(c.GetHeader("Authorization"))
		if !ok {
			abortUnauthorized(c, "token_missing",
				"Header Authorization tidak ada atau formatnya salah (harus: Bearer <token>)")
			return
		}

		claims := &jwt.RegisteredClaims{}

		// jwt.ParseWithClaims melakukan tiga hal sekaligus:
		//  1. decode payload token ke dalam `claims`
		//  2. memverifikasi signature memakai jwtKey
		//  3. memvalidasi waktu (exp / nbf) — otomatis di golang-jwt v5
		//
		// Di v5, kalau err == nil maka token PASTI valid, jadi pengecekan
		// `!token.Valid` yang terpisah sudah tidak diperlukan lagi.
		_, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (any, error) {
				// PENTING: pastikan algoritmanya HMAC.
				// Tanpa pengecekan ini, kode rentan terhadap serangan
				// "algorithm confusion" — penyerang mengirim token dengan
				// algoritma lain agar kunci kita diperlakukan sebagai
				// kunci publik.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return jwtKey, nil
			},
			// Lapis pertahanan kedua: hanya HS256 yang diterima.
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)

		if err != nil {
			// Token kedaluwarsa dibedakan dari token palsu supaya frontend
			// bisa memutuskan: refresh token, atau paksa login ulang.
			// Keduanya sama-sama 401, yang membedakan adalah field "code".
			if errors.Is(err, jwt.ErrTokenExpired) {
				abortUnauthorized(c, "token_expired",
					"Token sudah kedaluwarsa, silakan login ulang")
				return
			}
			abortUnauthorized(c, "token_invalid", "Token tidak valid")
			return
		}

		// Token bisa saja valid secara kriptografis tapi tidak memuat
		// identitas (Subject kosong). Kalau lolos, controller akan
		// menerima username kosong tanpa sadar.
		if claims.Subject == "" {
			abortUnauthorized(c, "token_invalid", "Token tidak memuat identitas user")
			return
		}

		c.Set(ContextUsername, claims.Subject)
		c.Next()
	}
}

// extractBearerToken memisahkan token dari header Authorization.
//
// Lebih ketat daripada strings.TrimPrefix yang dipakai versi lama:
//   - header tanpa kata "Bearer" ditolak, bukan diam-diam diproses
//   - "bearer" huruf kecil tetap diterima (EqualFold), karena beberapa
//     HTTP client menormalkan kapitalisasi
//   - token yang isinya hanya spasi dianggap kosong
func extractBearerToken(header string) (string, bool) {
	if len(header) <= len(bearerPrefix) {
		return "", false
	}
	if !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return "", false
	}

	token := strings.TrimSpace(header[len(bearerPrefix):])
	return token, token != ""
}
