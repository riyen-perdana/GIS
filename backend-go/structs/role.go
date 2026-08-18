package structs

/*
	TODO: Struct Role For Create
*/

type RoleCreateRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []uint `json:"permission_ids"`
}

/*
	TODO: Struct Role For Update
*/
type RoleUpdateRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []uint `json:"permission_ids"`
}

/*
	TODO: Struct Role For Show
*/
type RoleResponse struct {
	Id          uint                 `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}
