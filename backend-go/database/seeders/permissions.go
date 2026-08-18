package seeders

import (
	"backend-go/models"

	"gorm.io/gorm"
)

// SeedPermissions mengisi tabel permissions dengan daftar izin akses (permission)
// yang dibutuhkan aplikasi. Setiap permission mewakili satu aksi spesifik pada
// satu modul, dengan pola penamaan "modul-aksi" (contoh: "users-create" berarti
// izin untuk membuat data user).
//
// Permission-permission ini nantinya akan dihubungkan ke role (lihat seeder roles)
// sehingga akses user diatur berdasarkan role yang dimilikinya, bukan permission
// satu per satu.
func SeedPermissions(db *gorm.DB) {
	// Daftar seluruh permission yang tersedia di aplikasi, dikelompokkan per modul.
	permissions := []models.Permission{

		// Modul Dashboard: hanya butuh akses untuk melihat halaman utama.
		{Name: "dashboard-index"},

		// Modul Users (manajemen pengguna): CRUD lengkap.
		{Name: "users-index"},  // melihat daftar user
		{Name: "users-create"}, // menambah user baru
		{Name: "users-show"},   // melihat detail satu user
		{Name: "users-edit"},   // membuka form edit user
		{Name: "users-update"}, // menyimpan perubahan data user
		{Name: "users-delete"}, // menghapus user

		// Modul Permissions (manajemen izin akses): CRUD lengkap.
		{Name: "permissions-index"},
		{Name: "permissions-create"},
		{Name: "permissions-show"},
		{Name: "permissions-edit"},
		{Name: "permissions-update"},
		{Name: "permissions-delete"},

		// Modul Roles (manajemen peran/jabatan): CRUD lengkap.
		{Name: "roles-index"},
		{Name: "roles-create"},
		{Name: "roles-show"},
		{Name: "roles-edit"},
		{Name: "roles-update"},
		{Name: "roles-delete"},

		// Modul Categories (kategori data, misal kategori lokasi/peta): CRUD lengkap.
		{Name: "categories-index"},
		{Name: "categories-create"},
		{Name: "categories-show"},
		{Name: "categories-edit"},
		{Name: "categories-update"},
		{Name: "categories-delete"},

		// Modul Maps (data peta/lokasi): CRUD lengkap.
		{Name: "maps-index"},
		{Name: "maps-create"},
		{Name: "maps-show"},
		{Name: "maps-edit"},
		{Name: "maps-update"},
		{Name: "maps-delete"},

		// Modul Settings (pengaturan aplikasi): hanya perlu lihat & ubah,
		// tidak ada create/delete karena pengaturan biasanya berupa data tunggal.
		{Name: "settings-show"},
		{Name: "settings-update"},
	}

	// Simpan setiap permission ke database satu per satu.
	// FirstOrCreate akan mencari permission dengan Name yang sama terlebih dahulu;
	// jika belum ada, baru dibuat. Ini membuat seeder aman dijalankan berulang kali
	// (idempotent) tanpa menghasilkan data duplikat.
	for _, p := range permissions {
		db.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}
}
