package storage

import (
	"errors"
)

var (
	ErrorUserNotFound  = errors.New("user with given username doesn't exist")
	ErrorUserDuplicate = errors.New("user with given username already exists")

	ErrorMaxUrlCount = errors.New("user reached max url count")
)
