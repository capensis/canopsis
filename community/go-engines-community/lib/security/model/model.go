package model

const (
	ObjectTypeCRUD = "CRUD"
	ObjectTypeRW   = "RW"
)

// available permissions
const (
	PermissionCreate = "create"
	PermissionRead   = "read"
	PermissionUpdate = "update"
	PermissionDelete = "delete"
	PermissionCan    = "can"
)

// bitmasks of available permissions
const (
	PermissionBitmaskCreate = 8 // 0b1000
	PermissionBitmaskRead   = 4 // 0b0100
	PermissionBitmaskUpdate = 2 // 0b0010
	PermissionBitmaskDelete = 1 // 0b0001
	PermissionBitmaskCan    = 1 // 0b0001
)
