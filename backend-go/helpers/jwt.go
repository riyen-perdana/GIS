package helpers

import (
	"backend-go/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtKey di sini HARUS sama persis dengan jwtKey di middlewares/auth.go,
// karena keduanya memakai env var yang sama (JWT_KEY): jwt.go menandatangani
// token, auth.go memverifikasi tanda tangan itu. Kalau nilainya beda,
// token yang dibuat GenerateToken akan selalu ditolak sebagai "token_invalid".
var jwtKey = mustJWTKey()

// mustJWTKey mengambil JWT_KEY dari .env, dan menghentikan aplikasi (fatal)
// kalau belum diset — sengaja tidak diberi nilai fallback seperti
// "secret_key", supaya tidak ada yang lupa mengisi .env lalu tanpa sadar
// menandatangani token pakai kunci yang gampang ditebak.
func mustJWTKey() []byte {
	key := config.GetEnv("JWT_KEY", "")
	if key == "" {
		log.Fatal("[FATAL] JWT_KEY belum di-set. Tambahkan JWT_KEY di file .env sebelum menjalankan aplikasi.")
	}
	return []byte(key)
}

func GenerateToken(username string) string {

	// Mengatur waktu kedaluwarsa token, di sini kita set 60 menit dari waktu sekarang
	expirationTime := time.Now().Add(60 * time.Minute)

	// Membuat klaim (claims) JWT
	// Subject berisi username, dan ExpiresAt menentukan waktu expired token
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Membuat token baru dengan klaim yang telah dibuat
	// Menggunakan algoritma HS256 untuk menandatangani token
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	// Mengembalikan token dalam bentuk string
	return token
}
