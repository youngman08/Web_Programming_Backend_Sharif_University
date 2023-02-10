package main

import (
	"strconv"
	"time"

	"github.com/floydjones1/auth-server/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Name           string
	PassportNumber string
	Email          string
	Password       string
}

type LoginRequest struct {
	Email    string
	Password string
}

func main() {
	engine, err := data.CreateDBEngine()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Post("/signup", func(c *fiber.Ctx) error {
		req := new(SignupRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		if req.Name == "" || req.Email == "" || req.Password == "" || req.PassportNumber == "" {
			return c.Status(404).JSON(fiber.Map{"error": "invalid signup credentials"})
		}

		// save this info in the database
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		passNumber, err := strconv.Atoi(req.PassportNumber)

		if err != nil {
			return err
		}

		user := &data.User{
			Name:           req.Name,
			PassportNumber: passNumber,
			Email:          req.Email,
			Password:       string(hash),
		}

		_, err = engine.Insert(user)
		if err != nil {
			return err
		}
		token, exp, err := createJWTToken(*user)
		if err != nil {
			return err
		}
		// create a jwt token

		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		req := new(LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		if req.Email == "" || req.Password == "" {
			return c.Status(404).JSON(fiber.Map{"error": "invalid login credentials"})
		}

		user := new(data.User)
		has, err := engine.Where("email = ?", req.Email).Desc("id").Get(user)
		if err != nil {
			return err
		}
		if !has {
			return c.Status(404).JSON(fiber.Map{"error": "invalid login credentials"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return err
		}

		token, exp, err := createJWTToken(*user)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	private := app.Group("/private")
	private.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))
	private.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"success": true, "path": "private"})
	})

	public := app.Group("/public")
	public.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"success": true, "path": "public"})
	})

	if err := app.Listen("localhost:3001"); err != nil {
		panic(err)
	}
}

func createJWTToken(user data.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 60).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
