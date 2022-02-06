package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arvan-sistan/phonebook-rest-api/internal/http/request"
	"github.com/arvan-sistan/phonebook-rest-api/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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

func (u User) Login(c *fiber.Ctx) error {
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

	// loading user from db
	user, err := u.Storage.LoadByUserPass(c.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			return fiber.ErrNotFound
		}

		log.Printf("cannot load user %s", err)

		return fiber.ErrInternalServerError
	}

	// creating jwt claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // read expire from configmap
	}

	// signing jwt claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// read the secret from configmap
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": t})
}

// registering urls
func (u User) Register(g fiber.Router) {
	g.Post("/signup", u.SignUp)
	g.Post("/login", u.Login)
}
