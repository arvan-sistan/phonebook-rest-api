package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/arvan-sistan/phonebook-rest-api/internal/http/request"
	"github.com/arvan-sistan/phonebook-rest-api/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Storage storage.User
}

func (u User) SignUp(c *fiber.Ctx) error {

	req := new(request.User)

	// gets body of user's json request
	if err := c.BodyParser(req); err != nil {
		log.Printf("cannot load user data %s", err)

		return fiber.ErrBadRequest
	}

	//checks if user's request is valid
	if err := req.Validate(); err != nil {
		log.Printf("cannot validate user data %s", err)

		return fiber.ErrBadRequest
	}

	// saves user in db
	user, err := u.Storage.SaveUser(c.Context(), *req)
	if err != nil {
		if errors.Is(err, storage.ErrorUserDuplicate) {
			return fiber.NewError(http.StatusBadRequest, "user already exists")
		}

		log.Printf("cannot save user %s", err)

		return fiber.ErrInternalServerError
	}

	return c.Status(http.StatusCreated).JSON(user)

}

func (u User) Register(g fiber.Router) {
	g.Post("/signup", u.SignUp)
}
