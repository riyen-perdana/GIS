// package database: kumpulan fungsi bantu untuk mengelola koneksi ke database
// (bukan package main, hanya library yang dipanggil dari file lain, misal main.go)
package database

import (
	// fmt: untuk membuat string terformat (Sprintf) dan mencetak ke terminal (Println)
	"fmt"
	// log: untuk mencetak pesan error lalu menghentikan program (log.Fatal)
	"log"

	// driver mysql untuk GORM, dipakai gorm.Open() supaya tahu cara bicara ke database MySQL
	"gorm.io/driver/mysql"
	// gorm: ORM (Object-Relational Mapping) untuk Go, mempermudah interaksi ke database
	// tanpa perlu menulis query SQL manual satu-satu
	"gorm.io/gorm"

	// package config buatan sendiri, dipakai untuk baca environment variable (.env)
	"backend-go/config"
	// package models buatan sendiri, berisi struct-struct GORM (User, Role,
	// Permission, Category, Map, Setting) yang dipakai AutoMigrate di bawah
	"backend-go/models"
)

// DB adalah variabel global yang menyimpan koneksi database aktif.
// Setelah InitDb() dipanggil (biasanya sekali di awal, dari main.go),
// variabel ini bisa dipakai di file/package lain untuk query,
// misal: database.DB.Find(&users)
var DB *gorm.DB

// InitDb membuka koneksi ke database MySQL menggunakan GORM,
// lalu menyimpan hasilnya ke variabel global DB di atas.
// Kalau koneksi gagal, program langsung dihentikan (lewat log.Fatal),
// karena aplikasi backend biasanya tidak berguna tanpa database.
func InitDb() {

	// Ambil kredensial/konfigurasi database dari environment variable.
	// Kalau env variable-nya tidak di-set, pakai nilai default di sebelahnya
	// (memudahkan saat development lokal tanpa perlu setup .env dulu)
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "3306")
	dbUser := config.GetEnv("DB_USER", "root")
	dbPass := config.GetEnv("DB_PASS", "")
	dbName := config.GetEnv("DB_NAME", "db_golang_gis")

	// DSN (Data Source Name) adalah string berisi semua info yang dibutuhkan
	// untuk membuka koneksi: user, password, host, port, nama database, dst.
	// Formatnya spesifik mengikuti aturan driver mysql yang dipakai.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	// Membuat koneksi ke database
	// gorm.Open(driver, config) mencoba konek ke database sesuai DSN di atas.
	// Hasilnya 2 nilai: koneksi (*gorm.DB) dan error (nil kalau berhasil).
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Fatal mencetak pesan error lalu langsung menghentikan program (os.Exit(1)).
		// Dipakai di sini karena tanpa koneksi database, aplikasi tidak bisa jalan normal,
		// jadi lebih baik program berhenti dari awal daripada lanjut dalam keadaan rusak.
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Auto Migrate Models
	// AutoMigrate membuat tabel yang belum ada dan menyesuaikan skema tabel
	// yang sudah ada (kolom, index, foreign key) supaya cocok dengan struct
	// model di package models — tanpa perlu menulis migration SQL manual.
	// GORM otomatis mengurutkan pembuatan tabel sesuai dependensi foreign key
	// (mis. Category dibuat sebelum Map), jadi urutan argumen di bawah bebas.
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Category{}, &models.Map{}, &models.Setting{}) // Tambahkan model lain jika perlu
	if err != nil {
		// Migrasi gagal (mis. konflik skema/constraint) — hentikan program
		// karena tabel yang tidak sinkron dengan model bisa menyebabkan bug
		// yang lebih sulit dilacak kalau aplikasi tetap dilanjutkan.
		log.Fatal("Gagal Melakukan Migrasi Database:", err)
	}

	fmt.Println("Koneksi ke database berhasil")

}
