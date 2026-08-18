package structs

/*
	TODO: Struct User For Create Data
*/
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required"`
	RoleIDs  []uint `json:"role_ids" binding:"required"`
}

/*
	TODO: Struct User For Update Data
*/
type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required"`
	RoleIDs  []uint `json:"role_ids" binding:"required"`
}

/*
	TODO: Struct User For Show Data
*/
type UserResponse struct {
	Id          uint                 `json:"id"`
	Name        string               `json:"name"`
	Username    string               `json:"username"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	Roles       []RoleResponse       `json:"roles"`
	Email       string               `json:"email"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type UserSimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
