package middlewares

import (
	"errors"
	"net/http"

	"backend-go/database"
	"backend-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Kunci-kunci yang dipakai untuk menyimpan data di gin.Context.
// Dibuat konstanta supaya tidak ada typo "usernme" yang baru ketahuan saat runtime.
const (
	ContextUsername    = "username"
	ContextCurrentUser = "currentUser"
)

// ErrNoUsernameInContext muncul kalau Permission() dipasang di route
// yang lupa memakai Auth() lebih dulu.
var ErrNoUsernameInContext = errors.New("username tidak ada di context")

// Username mengambil username hasil verifikasi token.
// c.GetString() sudah melakukan type assertion dan mengembalikan ""
// kalau key tidak ada — jadi lebih aman daripada c.Get() yang bertipe any.
func Username(c *gin.Context) string {
	return c.GetString(ContextUsername)
}

// CurrentUser mengembalikan user yang sedang login, LENGKAP dengan
// relasi Roles dan Permissions.
//
// Hasil query disimpan ke context, jadi kalau dipanggil berkali-kali
// dalam satu request (oleh middleware Permission, lalu oleh controller),
// database hanya disentuh SEKALI.
//
// Controller bisa memakainya juga:
//
//	user, err := middlewares.CurrentUser(c)
func CurrentUser(c *gin.Context) (*models.User, error) {
	if cached, ok := c.Get(ContextCurrentUser); ok {
		if user, ok := cached.(*models.User); ok {
			return user, nil
		}
	}

	username := Username(c)
	if username == "" {
		return nil, ErrNoUsernameInContext
	}

	var user models.User
	if err := database.DB.
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		// Error dikembalikan apa adanya supaya pemanggil bisa membedakan
		// "user tidak ada" (gorm.ErrRecordNotFound) dari "database bermasalah".
		return nil, err
	}

	c.Set(ContextCurrentUser, &user)
	return &user, nil
}

// abortUnauthorized: identitas tidak jelas → 401. Frontend sebaiknya
// mengarahkan user ke halaman login / melakukan refresh token.
//
// AbortWithStatusJSON = c.JSON(...) + c.Abort() dalam satu panggilan,
// sehingga tidak ada risiko lupa memanggil Abort.
// Pemanggil TETAP wajib menulis `return` setelahnya.
func abortUnauthorized(c *gin.Context, code, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": message,
		"code":  code,
	})
}

// abortForbidden: identitas sudah jelas dan valid, tapi haknya kurang → 403.
// Frontend TIDAK perlu menyuruh login ulang.
func abortForbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error": message,
		"code":  "permission_denied",
	})
}

// abortServerError: kesalahan di sisi kita (mis. koneksi DB putus) → 500.
// Penting: jangan mengembalikan 401 untuk error database, karena itu
// menyesatkan frontend seolah-olah token user yang bermasalah.
func abortServerError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": message,
		"code":  "internal_error",
	})
}

// isRecordNotFound dipakai untuk membedakan "data tidak ditemukan"
// dari error database lainnya.
func isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
