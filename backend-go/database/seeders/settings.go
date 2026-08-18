package seeders

import (
	// Package internal project, berisi struct model tabel database (Setting, dst).
	"backend-go/models"
	// "log" adalah package standar Go untuk mencatat pesan/error ke output
	// (biasanya terminal). Dipakai di sini untuk menghentikan program dengan pesan
	// yang jelas kalau proses seeding gagal.
	"log"

	// ORM untuk berkomunikasi dengan database lewat method Go, bukan SQL mentah.
	"gorm.io/gorm"
)

// SeedSetting mengisi SATU baris data pengaturan (setting) aplikasi, misalnya
// judul situs, titik tengah peta, dan level zoom default.
// Berbeda dari SeedRoles/SeedUsers yang meng-handle banyak data dalam loop,
// di sini kita cuma butuh satu baris "konfigurasi" saja — makanya polanya lebih sederhana.
func SeedSetting(db *gorm.DB) {
	// `defaults` adalah struct literal berisi nilai-nilai default yang akan dipakai
	// HANYA JIKA data setting belum ada di database (lihat penjelasan Attrs di bawah).
	// Nilai MapCenterLat/Lng disimpan sebagai string (bukan float) karena tipe kolom
	// di database adalah "double" tapi field Go-nya string — kemungkinan untuk
	// menghindari masalah pembulatan angka desimal saat marshal/unmarshal JSON.
	// VillageBoundary diisi "[]" (string JSON array kosong) sebagai placeholder,
	// nantinya bisa diisi data GeoJSON batas wilayah desa.
	defaults := models.Setting{
		Title:           "GIS Desa Santri",
		Description:     "Eksplorasi Desa Santri secara interaktif melalui peta GIS.",
		MapCenterLat:    "-7.592589928951457",
		MapCenterLng:    "112.26113954274147",
		MapZoom:         16,
		VillageBoundary: "[]",
	}

	// `var out models.Setting` menyiapkan variabel kosong (zero value) untuk
	// menampung hasil query di bawah — baik itu data yang sudah ada, maupun
	// data baru yang baru saja dibuat.
	var out models.Setting
	// Blok ini memakai "method chaining" ala GORM: beberapa method dipanggil
	// berurutan (Where -> Attrs -> FirstOrCreate), masing-masing mengembalikan
	// *gorm.DB yang sama supaya bisa disambung lagi, sampai akhirnya kita
	// baca field `.Error` di ujung chain.
	if err := db.
		// Where(&models.Setting{Id: 1}): kondisi PENCARIAN — cari baris setting
		// dengan Id = 1. Karena aplikasi ini cuma butuh satu baris pengaturan,
		// Id 1 dipakai sebagai "kunci tunggal" (mirip singleton row).
		Where(&models.Setting{Id: 1}).
		// Attrs(defaults): berbeda dari Where, Attrs TIDAK dipakai untuk mencari,
		// tapi disiapkan sebagai nilai yang akan dipakai HANYA saat data belum
		// ditemukan dan GORM perlu membuat baris baru. Jadi kombinasi
		// Where + Attrs berarti: "cari Id=1, kalau ketemu pakai apa adanya,
		// kalau tidak ketemu buat baru dengan Id=1 + isi dari `defaults`".
		Attrs(defaults).
		// FirstOrCreate(&out): jalankan pencarian di atas. Hasilnya (data lama
		// atau data baru yang baru dibuat) ditulis ke variabel `out` lewat pointer.
		// `.Error` di akhir mengambil error terakhir dari keseluruhan chain ini
		// (langsung diakses tanpa variabel perantara, gaya ringkas ala Go).
		FirstOrCreate(&out).Error; err != nil {
		// `if err := ...; err != nil` adalah idiom Go: mendeklarasikan variabel
		// `err` SEKALIGUS langsung mengeceknya dalam satu baris `if`. Variabel
		// `err` ini hanya berlaku (scope-nya) di dalam blok if ini saja.
		// log.Fatalf mencetak pesan error lalu langsung menghentikan program
		// (memanggil os.Exit(1) di baliknya) — dipakai karena kalau setting
		// dasar gagal disiapkan, aplikasi tidak masuk akal untuk lanjut jalan.
		log.Fatalf("failed seeding setting: %v", err)
	}
}
