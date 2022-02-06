package storage

import (
	"context"
	"errors"

	"github.com/arvan-sistan/phonebook-rest-api/internal/http/request"
	"github.com/arvan-sistan/phonebook-rest-api/internal/model"
)

var (
	ErrorUserNotFound  = errors.New("user with given username doesn't exist")
	ErrorUserDuplicate = errors.New("user with given username already exists")
)

type User interface {
	SaveUser(context.Context, request.User) (model.User, error)
	LoadByUserPass(context.Context, string, string) (model.User, error)
}
