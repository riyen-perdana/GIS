package structs

/*
	TODO: Struct Permission For Create Data
*/
type PermissionCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

/*
	TODO: Struct Permission For Update Data
*/
type PermissionUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

/*
	TODO: Struct Permission For Show Data
*/
type PermissionResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
