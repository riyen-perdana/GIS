package seeders

import (
	// Package internal project yang menyimpan variabel global `database.DB`
	// (koneksi *gorm.DB yang sudah diinisialisasi lewat InitDb, lihat database/database.go).
	"backend-go/database"
	// Package standar Go untuk mencetak log ke terminal — dipakai di sini
	// hanya untuk menampilkan progress, bukan untuk menghentikan program.
	"log"
)

// Seed adalah "pintu masuk" (entry point) untuk menjalankan semua seeder sekaligus.
// Fungsi ini biasanya dipanggil dari file terpisah, misalnya sebuah command/cmd
// khusus (contoh: `go run cmd/seed/main.go`), bukan dari server utama —
// tujuannya supaya proses isi data awal ini terpisah dari alur jalannya aplikasi.
func Seed() {
	// Ambil koneksi database yang sudah dibuat sebelumnya (variabel global).
	// Karena `database.DB` bertipe *gorm.DB (pointer), variabel lokal `db` di sini
	// menunjuk ke koneksi yang SAMA, bukan salinannya.
	db := database.DB
	// log.Println mencetak pesan biasa (bukan error) — berguna untuk memantau
	// proses seeding di terminal, terutama kalau datanya banyak dan makan waktu.
	log.Println("Running database seeders...")

	// Jalankan seeder secara berurutan.
	// Urutan pemanggilan di sini PENTING karena ada ketergantungan data:
	// - SeedPermissions harus jalan duluan, karena SeedRoles butuh daftar
	//   permission yang sudah ada untuk dikaitkan ke role (lihat roles.go).
	// - SeedRoles harus jalan sebelum SeedUsers, karena SeedUsers mencari
	//   role "admin"/"user" yang sudah tersimpan untuk dikaitkan ke user (lihat users.go).
	// - SeedSetting tidak bergantung pada seeder lain, jadi posisinya bebas,
	//   tapi diletakkan terakhir sebagai konvensi "data pelengkap".
	// Setiap fungsi seeder dipanggil secara SEKUENSIAL (satu selesai baru lanjut
	// ke berikutnya) — bukan konkuren/goroutine — supaya urutan ketergantungan
	// data di atas tetap terjaga.
	SeedPermissions(db)
	SeedRoles(db)
	SeedUsers(db)
	SeedSetting(db)

	// Ditampilkan setelah semua seeder di atas selesai dijalankan tanpa panic/Fatalf.
	log.Println("Database seeding completed!")
}
