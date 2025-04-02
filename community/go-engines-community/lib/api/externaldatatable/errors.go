package externaldatatable

import "errors"

var ErrConfigNotDeletable = errors.New("table from config cannot be deleted")
var ErrLinkedNotDeletable = errors.New("linked table cannot be deleted")
