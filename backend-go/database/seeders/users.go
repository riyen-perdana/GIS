package seeders

import (
	// Package internal project, berisi struct model tabel database (User, Role, dst).
	"backend-go/models"

	// golang.org/x/crypto/bcrypt: library standar-de-facto Go untuk hashing password.
	// Password TIDAK BOLEH disimpan sebagai plain text di database — bcrypt mengubahnya
	// jadi hash satu arah (tidak bisa dibalik/di-decrypt) yang aman untuk disimpan.
	"golang.org/x/crypto/bcrypt"
	// ORM untuk berkomunikasi dengan database lewat method Go, bukan SQL mentah.
	"gorm.io/gorm"
)

// SeedUsers mengisi data user awal (admin & user biasa) ke database.
// Ditulis idempotent (aman dijalankan berkali-kali): kalau user sudah ada, datanya
// akan di-update, bukan membuat duplikat.
func SeedUsers(db *gorm.DB) {
	// bcrypt.GenerateFromPassword mengembalikan DUA nilai: hasil hash ([]byte) dan error.
	// Di Go, hampir semua operasi yang bisa gagal mengembalikan error sebagai nilai balik
	// terakhir — bukan lewat exception seperti di bahasa lain (Java/Python/JS).
	// Di sini error-nya sengaja diabaikan dengan `_` karena bcrypt.DefaultCost adalah
	// konstanta valid bawaan library, jadi praktis tidak mungkin gagal.
	// bcrypt.DefaultCost adalah angka "cost factor" — makin tinggi, makin lambat & aman
	// proses hashing-nya (melindungi dari brute-force).
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	// `var adminRole, userRole models.Role` mendeklarasikan dua variabel struct Role
	// sekaligus, masing-masing diisi dengan zero value (semua field kosong/default).
	var adminRole, userRole models.Role
	// db.Where("name = ?", "admin").First(&adminRole):
	// - Where(...) menyusun klausa SQL "WHERE name = 'admin'" (tanda `?` diganti otomatis
	//   dan aman dari SQL injection oleh GORM).
	// - First(&adminRole) mengambil SATU baris pertama yang cocok, lalu menuliskan hasilnya
	//   ke variabel adminRole lewat pointer-nya.
	// Query ini mengasumsikan seeder role (lihat roles.go) sudah dijalankan lebih dulu,
	// karena kita butuh Id role yang sudah tersimpan di database.
	db.Where("name = ?", "admin").First(&adminRole)
	db.Where("name = ?", "user").First(&userRole)

	// Slice berisi data user yang ingin kita pastikan ada di database.
	// Field Roles diisi dengan slice berisi satu Role (relasi many2many, lihat tag
	// `gorm:"many2many:user_roles"` di models/user.go) — artinya satu user bisa
	// punya banyak role, di sini kita hanya kaitkan satu.
	users := []models.User{
		{
			Name:     "Admin",
			Username: "admin",
			Email:    "admin@gmail.com",
			Password: string(password), // []byte hasil hash diubah jadi string agar cocok dengan tipe field Password.
			Roles:    []models.Role{adminRole},
		},
		{
			Name:     "User",
			Username: "user",
			Email:    "user@gmail.com",
			Password: string(password),
			Roles:    []models.Role{userRole},
		},
	}

	// Iterasi setiap user yang ingin diseed. `u` adalah salinan (copy) dari tiap elemen.
	for _, u := range users {
		var user models.User
		// Cek dulu apakah user dengan username ini SUDAH ada di database.
		// `.Error` adalah field pada hasil query GORM (*gorm.DB) yang menyimpan
		// error terakhir yang terjadi — pola umum GORM: chain method dulu,
		// baru cek `.Error` di akhir.
		err := db.Where("username = ?", u.Username).First(&user).Error

		// Pola idiomatis Go: `if err != nil { ... }` untuk menangani kegagalan.
		// Berbeda dari try-catch, di sini kita HARUS mengecek error secara eksplisit
		// setiap kali sebuah operasi bisa gagal.
		if err != nil {
			// gorm.ErrRecordNotFound adalah error khusus bawaan GORM yang berarti
			// "query berhasil dijalankan, tapi tidak ada baris yang cocok" —
			// ini BUKAN error sungguhan, melainkan sinyal normal "data belum ada".
			if err == gorm.ErrRecordNotFound {
				// Data belum ada, buat baru.
				// db.Create(&u) akan INSERT user baru sekaligus otomatis membuat
				// relasi many2many ke Roles (lewat tabel pivot "user_roles"),
				// karena GORM mendukung "auto-save associations" saat Create.
				db.Create(&u)
			} else {
				// Kalau error-nya BUKAN "record not found", berarti ada masalah
				// yang lebih serius (misalnya koneksi database putus).
				// `panic(err)` menghentikan program secara paksa — dipakai di sini
				// karena seeder adalah proses sekali-jalan (bukan HTTP handler),
				// jadi kalau gagal total, lebih baik langsung berhenti dan
				// menampilkan error-nya daripada melanjutkan dengan data yang salah.
				panic(err)
			}
		} else {
			// User SUDAH ada (tidak ada error) → update saja datanya, jangan duplikat.
			// db.Model(&user).Updates(...): meng-update baris `user` yang sudah
			// ditemukan di atas, hanya field yang diisi di struct models.User{...}
			// (Email dan Password) yang akan diubah — field lain (Name, dll) tidak disentuh.
			db.Model(&user).Updates(models.User{
				Email:    u.Email,
				Password: string(password),
			})
			// Perbarui juga daftar Roles user ini, mengganti relasi lama dengan
			// relasi baru sesuai data seeder (u.Roles).
			db.Model(&user).Association("Roles").Replace(u.Roles)
		}
	}

}
