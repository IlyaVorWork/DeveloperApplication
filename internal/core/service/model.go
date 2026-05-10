package service

import "errors"

var ErrNotVerified = errors.New("application must be successfully verified before publishing")
