// package config: kumpulan fungsi bantu untuk mengelola konfigurasi aplikasi
// (bukan package main, karena file ini bukan entry point, hanya "library"
// yang dipakai/dipanggil dari file lain, misal dari main.go)
package config

import (
	// log: untuk mencetak pesan ke terminal, khusus untuk log/peringatan
	"log"
	// os: untuk mengakses environment variable milik sistem operasi
	"os"

	// godotenv: library pihak ketiga untuk membaca file .env
	// dan memuat isinya sebagai environment variable
	"github.com/joho/godotenv"
)

// LoadEnv membaca file .env di root project (kalau ada) lalu memuat
// isinya sebagai environment variable, supaya bisa diambil dengan os.Getenv/GetEnv.
// Biasanya dipanggil sekali di awal, di dalam func main().
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		// Kalau file .env tidak ditemukan, program tidak dihentikan (tidak pakai panic/os.Exit).
		// Cuma dikasih peringatan, karena env variable bisa saja sudah di-set
		// dengan cara lain, misal langsung di sistem/server (tanpa file .env).
		log.Println("Warning : Tidak ada file .env yang ditemukan, menggunakan variabel lingkungan default")
	}
}

// GetEnv mengambil nilai environment variable berdasarkan key (nama variabelnya).
// Kalau variabel dengan nama tersebut tidak ditemukan, akan mengembalikan
// defaultValue supaya program tidak error/crash gara-gara env variable kosong.
//
// Contoh pemakaian:
//
//	port := config.GetEnv("PORT", "3000") // kalau PORT tidak di-set, pakai "3000"
func GetEnv(key string, defaultValue string) string {
	// os.LookupEnv mengembalikan 2 nilai:
	// value  -> isi variabelnya (kalau ada)
	// exists -> true/false, apakah variabel itu benar-benar di-set atau tidak
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
