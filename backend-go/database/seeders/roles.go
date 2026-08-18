// Package seeders berisi kumpulan fungsi untuk mengisi ("seed") data awal ke database,
// biasanya dipanggil sekali saat setup aplikasi atau lewat perintah `go run` khusus seeder.
package seeders

import (
	// "backend-go/models" adalah package lokal (internal) milik project ini,
	// berisi definisi struct yang merepresentasikan tabel-tabel database (Role, Permission, dst).
	// Import path-nya mengikuti nama module di go.mod (lihat baris `module backend-go`).
	"backend-go/models"

	// gorm.io/gorm adalah library ORM (Object Relational Mapping) untuk Go.
	// ORM memungkinkan kita menulis query database (SELECT, INSERT, UPDATE, dll)
	// menggunakan method Go biasa, tanpa menulis SQL mentah secara manual.
	"gorm.io/gorm"
)

// SeedRoles adalah fungsi biasa (bukan method, karena tidak ada receiver sebelum nama fungsi).
// Parameter `db *gorm.DB` adalah pointer ke koneksi database yang sudah terhubung —
// pointer dipakai supaya fungsi ini memakai koneksi yang sama persis dengan yang dibuat
// di tempat lain (misalnya main.go), bukan salinan/copy-nya.
// Fungsi ini tidak mengembalikan nilai apa pun (tidak ada return type), jadi tugasnya
// murni "efek samping": menulis data ke database.
func SeedRoles(db *gorm.DB) {

	// `roles` adalah slice (mirip array tapi ukurannya dinamis) berisi 2 buah
	// struct literal bertipe models.Role. Di sini kita hanya mengisi field `Name`,
	// field lain (Id, Permissions, CreatedAt, dst) akan diisi otomatis dengan nilai
	// default/zero value oleh Go (misalnya Id = 0, CreatedAt = time.Time kosong).
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	// `for _, role := range roles` adalah cara Go melakukan iterasi (looping) atas slice.
	// `range` mengembalikan dua nilai: index dan value. Karena kita tidak butuh index-nya,
	// kita "buang" dengan underscore `_` (konvensi Go untuk mengabaikan suatu nilai).
	// `role` di sini adalah SALINAN (copy) dari tiap elemen roles, bukan referensi langsung.
	for _, role := range roles {
		// FirstOrCreate: perintah GORM yang artinya "cari data yang cocok dengan kondisi,
		// kalau tidak ketemu, buat baru". Di sini kondisinya adalah models.Role{Name: role.Name}
		// (mencari role dengan nama yang sama), dan hasilnya (baik data lama maupun yang baru
		// dibuat) akan dimasukkan kembali ke variabel `role` lewat pointer `&role`.
		// Tujuannya: seeder ini aman dijalankan berkali-kali (idempotent) — tidak akan
		// membuat data duplikat kalau role "admin"/"user" sudah ada sebelumnya.
		db.FirstOrCreate(&role, models.Role{Name: role.Name})

		// `var allPermissions []models.Permission` mendeklarasikan slice kosong (nil slice)
		// bertipe Permission. Ini akan diisi oleh query di bawah.
		var allPermissions []models.Permission
		// db.Find(&allPermissions): mengambil SEMUA baris dari tabel permissions
		// (tanpa filter/Where), lalu hasilnya dimasukkan ke slice allPermissions
		// lewat pointer-nya (&allPermissions), karena GORM perlu mengubah isi
		// variabel tersebut secara langsung.
		db.Find(&allPermissions)

		// `switch role.Name { ... }` adalah percabangan mirip if-else berantai,
		// tapi lebih rapi ketika membandingkan satu variabel dengan banyak kemungkinan nilai.
		// Berbeda dari bahasa lain, di Go setiap `case` otomatis "break" (tidak perlu
		// menulis break manual, dan tidak akan lanjut/fallthrough ke case berikutnya).
		switch role.Name {
		case "admin":
			// Untuk role "admin", kita berikan SEMUA permission yang ada.
			// db.Model(&role) memberi tahu GORM "kita sedang bekerja pada record role ini".
			// .Association("Permissions") mengacu ke relasi many2many yang didefinisikan
			// di struct Role (lihat tag `gorm:"many2many:role_permissions"` di models/role.go) —
			// GORM otomatis mengelola tabel pivot/junction "role_permissions" di belakang layar.
			// .Replace(allPermissions) akan MENGGANTI seluruh daftar permission role ini
			// dengan isi allPermissions (menghapus relasi lama, memasang relasi baru).
			db.Model(&role).Association("Permissions").Replace(allPermissions)
		case "user":
			// Untuk role "user", kita hanya berikan sebagian permission (akses terbatas/view-only).
			var viewOnly []models.Permission
			// db.Where("name IN ?", []string{...}): membangun query SQL
			// `WHERE name IN (...)` secara aman (parameter di-escape otomatis oleh GORM,
			// jadi tidak rawan SQL injection walau nilainya datang dari slice Go).
			// Daftar string ini adalah nama-nama permission spesifik yang boleh diakses user biasa,
			// misalnya untuk halaman kategori (categories) dan peta (maps): index (lihat list),
			// create (tambah), show (lihat detail), dan edit (ubah) — tanpa hak "delete".
			db.Where("name IN ?", []string{"categories-index", "categories-create", "categories-show", "categories-edit", "maps-index", "maps-create", "maps-show", "maps-edit"}).Find(&viewOnly)
			// Sama seperti admin, tapi menggunakan viewOnly (daftar permission terbatas)
			// alih-alih allPermissions.
			db.Model(&role).Association("Permissions").Replace(viewOnly)
		}
	}
}
