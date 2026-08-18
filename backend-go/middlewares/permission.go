package middlewares

import (
	"errors"

	"backend-go/models"

	"github.com/gin-gonic/gin"
)

// Permission mengizinkan request kalau user punya SALAH SATU dari
// permission yang disebutkan (logika OR).
//
//	middlewares.Permission("user.edit")                  // butuh satu permission
//	middlewares.Permission("user.edit", "user.manage")   // cukup punya salah satu
//
// Middleware ini mengasumsikan Auth() sudah berjalan lebih dulu.
// Kalau tidak, request ditolak dengan 401 — bukan 403 — karena masalahnya
// memang "kita tidak tahu kamu siapa", bukan "kamu tidak berhak".
func Permission(permissionNames ...string) gin.HandlerFunc {
	// Validasi ini berjalan saat route DIDAFTARKAN (waktu startup),
	// bukan saat ada request. Jadi kesalahan pemakaian langsung ketahuan
	// begitu aplikasi dijalankan — bukan diam-diam meloloskan semua orang.
	if len(permissionNames) == 0 {
		panic("middlewares.Permission: minimal butuh satu nama permission")
	}

	// Nama permission diubah jadi map sekali saja di luar closure.
	// Pencarian jadi O(1) dan tidak diulang tiap request.
	// struct{}{} dipakai karena tidak memakan memori — kita hanya
	// butuh keberadaan key-nya, bukan value.
	wanted := make(map[string]struct{}, len(permissionNames))
	for _, name := range permissionNames {
		wanted[name] = struct{}{}
	}

	return func(c *gin.Context) {
		user, err := CurrentUser(c)
		if err != nil {
			handleUserLoadError(c, err)
			return
		}

		if userHasAnyPermission(user, wanted) {
			c.Next()
			return
		}

		abortForbidden(c, "Anda tidak memiliki hak akses untuk aksi ini")
	}
}

// PermissionAll mewajibkan user punya SEMUA permission yang disebutkan
// (logika AND). Berguna untuk aksi sensitif yang menggabungkan dua wewenang.
//
//	middlewares.PermissionAll("data.export", "data.view")
func PermissionAll(permissionNames ...string) gin.HandlerFunc {
	if len(permissionNames) == 0 {
		panic("middlewares.PermissionAll: minimal butuh satu nama permission")
	}

	required := make(map[string]struct{}, len(permissionNames))
	for _, name := range permissionNames {
		required[name] = struct{}{}
	}

	return func(c *gin.Context) {
		user, err := CurrentUser(c)
		if err != nil {
			handleUserLoadError(c, err)
			return
		}

		owned := collectPermissions(user)
		for name := range required {
			if _, ok := owned[name]; !ok {
				abortForbidden(c, "Anda tidak memiliki hak akses untuk aksi ini")
				return
			}
		}

		c.Next()
	}
}

// Role mengizinkan request kalau user punya salah satu role tertentu.
// Pakai ini hanya untuk kasus yang benar-benar berbasis peran
// (mis. halaman khusus "admin"). Untuk aksi spesifik, Permission lebih
// baik karena tidak perlu mengubah kode saat pembagian peran berubah.
func Role(roleNames ...string) gin.HandlerFunc {
	if len(roleNames) == 0 {
		panic("middlewares.Role: minimal butuh satu nama role")
	}

	wanted := make(map[string]struct{}, len(roleNames))
	for _, name := range roleNames {
		wanted[name] = struct{}{}
	}

	return func(c *gin.Context) {
		user, err := CurrentUser(c)
		if err != nil {
			handleUserLoadError(c, err)
			return
		}

		for _, role := range user.Roles {
			if _, ok := wanted[role.Name]; ok {
				c.Next()
				return
			}
		}

		abortForbidden(c, "Role Anda tidak diizinkan mengakses resource ini")
	}
}

// handleUserLoadError memetakan error dari CurrentUser ke status HTTP
// yang tepat. Ini perbaikan penting dari versi lama yang mengembalikan
// 401 untuk SEMUA error — termasuk saat database sedang down, yang
// membuat frontend salah paham dan memaksa user login ulang percuma.
func handleUserLoadError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNoUsernameInContext):
		// Biasanya bug kita sendiri: route memakai Permission()
		// tapi lupa memasang Auth().
		abortUnauthorized(c, "unauthenticated", "Silakan login terlebih dahulu")
	case isRecordNotFound(err):
		// Token valid tapi user sudah dihapus/dinonaktifkan.
		abortUnauthorized(c, "user_not_found", "User pada token tidak ditemukan")
	default:
		abortServerError(c, "Gagal memuat data user")
	}
}

// userHasAnyPermission berhenti di kecocokan pertama, jadi tidak perlu
// menelusuri seluruh role kalau sudah ketemu.
func userHasAnyPermission(user *models.User, wanted map[string]struct{}) bool {
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if _, ok := wanted[perm.Name]; ok {
				return true
			}
		}
	}
	return false
}

// collectPermissions mengumpulkan seluruh nama permission milik user
// (gabungan dari semua role, otomatis terdeduplikasi oleh map).
func collectPermissions(user *models.User) map[string]struct{} {
	owned := make(map[string]struct{})
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			owned[perm.Name] = struct{}{}
		}
	}
	return owned
}
