package redis

import "errors"

var ErrFailedToRefreshLock = errors.New("failed to refresh lock")
