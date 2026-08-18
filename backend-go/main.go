// package main menandakan file ini adalah entry point program yang bisa dijalankan/dieksekusi
// (bukan sekadar library/package biasa yang cuma bisa dipakai oleh file lain)
package main

// import: memasukkan package yang dipakai di file ini
// - "backend-go/config" -> package lokal untuk memuat & membaca konfigurasi/env
// - "backend-go/database" -> package lokal untuk inisialisasi koneksi database
// - "backend-go/database/seeders" -> package lokal berisi fungsi Seed() untuk mengisi
//   data awal (role, user, setting, dst) ke database, lihat database/seeders/seeds.go
// - "github.com/gin-gonic/gin" -> package pihak ketiga, framework untuk bikin HTTP server
import (
	"backend-go/config"
	"backend-go/database"
	"backend-go/database/seeders"

	"github.com/gin-gonic/gin"
)

// func main() adalah fungsi yang otomatis dijalankan pertama kali
// saat program dieksekusi (go run main.go)
func main() {

	// Muat file .env (kalau ada) supaya nilai seperti APP_PORT
	// bisa diambil lewat config.GetEnv(). Wajib dipanggil di awal,
	// sebelum ada kode lain yang butuh baca environment variable.
	config.LoadEnv()

	// Buka koneksi ke database (lihat database/database.go).
	// Dipanggil setelah LoadEnv() karena InitDb() butuh baca env
	// seperti DB_HOST, DB_USER, dst yang baru tersedia setelah .env dimuat.
	database.InitDb()

	// Jalankan semua seeder (lihat database/seeders/seeds.go) setiap kali
	// aplikasi start. Dipanggil setelah InitDb() karena seeder butuh koneksi
	// database yang sudah aktif. Seeder-seedernya ditulis idempotent (aman
	// dijalankan berulang, tidak membuat data duplikat), jadi tidak masalah
	// kalau fungsi ini jalan lagi tiap kali server di-restart.
	seeders.Seed()

	// Inisialisasi Gin
	// gin.Default() membuat "router" baru lengkap dengan middleware bawaan
	// (logger untuk mencatat request masuk, dan recovery untuk menangkap panic/error)
	router := gin.Default()

	// Routing
	// router.GET(path, handler) mendaftarkan endpoint untuk method HTTP GET
	// "/" berarti endpoint ini aktif saat diakses di alamat root, mis. http://localhost:3000/
	router.GET("/", func(c *gin.Context) {
		// c (context) mewakili satu siklus request-response yang sedang berjalan
		// c.JSON(statusCode, data) mengirim balasan berupa JSON ke client
		// 200 = kode status HTTP untuk "OK" / berhasil
		// gin.H{...} adalah shortcut untuk membuat map[string]interface{} (mirip objek JSON)
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Menjalankan server
	// config.GetEnv("APP_PORT", "3000") -> ambil port dari environment variable APP_PORT,
	// kalau tidak di-set (mis. tidak ada file .env), pakai default "3000"
	// router.Run(":"+port) membuat server HTTP mendengarkan (listen) di port tersebut
	// setelah baris ini dijalankan, program akan terus berjalan menunggu request masuk
	// sampai dihentikan manual (mis. Ctrl+C di terminal)
	router.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
