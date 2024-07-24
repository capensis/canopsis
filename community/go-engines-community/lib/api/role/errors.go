package role

import "errors"

var ErrLinkedToUser = errors.New("role is linked to user")
var ErrUpdateAdminRole = errors.New("admin cannot be updated")
var ErrDeleteAdminRole = errors.New("admin cannot be deleted")
