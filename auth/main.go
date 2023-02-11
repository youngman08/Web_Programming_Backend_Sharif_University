package main

import (
	"strconv"
	"time"

	"github.com/floydjones1/auth-server/data"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisClient.Del("unauthorized")

	// Add db data
	var records []data.Unauthorized_token
	engine.Find(&records)
	for _, b := range records {
		redisClient.LPush("unauthorized", b)
	}

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
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
			FirstName:      req.Name,
			PassportNumber: passNumber,
			Email:          req.Email,
			PasswordHash:   string(hash),
		}

		result := engine.Create(user)
		if result.Error != nil {
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
		result := engine.First(&user, "email = ?", req.Email)
		if result.Error != nil {
			return c.Status(404).JSON(fiber.Map{"error": "invalid login credentials"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			return err
		}

		token, exp, err := createJWTToken(*user)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))
	app.Get("/info", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)

		for _, b := range redisClient.LRange("unauthorized", 0, -1).Val() {
			if b == user.Raw {
				return c.Status(400).JSON(fiber.Map{"Error": "Expired Token"})
			}
		}

		claims := user.Claims.(jwt.MapClaims)

		exp_time := time.Unix(int64(claims["exp"].(float64)), 0)
		now := time.Now()

		if exp_time.Sub(now) < 0 {
			// Add to expired tokens
			redisClient.LPush("unauthorized", user.Raw)

			token := &data.Unauthorized_token{
				UserId:     int64(claims["user_id"].(float64)),
				Token:      user.Raw,
				Expiration: exp_time,
			}
			result := engine.Create(token)
			if result.Error != nil {
				return err
			}
			return c.Status(400).JSON(fiber.Map{"Error": "Expired Token"})
		}

		return c.JSON(fiber.Map{"user_id": claims["user_id"], "PassportNumber": claims["PassportNumber"], "Name": claims["FirstName"]})
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)

		claims := user.Claims.(jwt.MapClaims)

		exp_time := time.Unix(int64(claims["exp"].(float64)), 0)

		redisClient.LPush("unauthorized", user.Raw)

		token := &data.Unauthorized_token{
			UserId:     int64(claims["user_id"].(float64)),
			Token:      user.Raw,
			Expiration: exp_time,
		}
		result := engine.Create(token)
		if result.Error != nil {
			return c.JSON(fiber.Map{"Logout": "Error"})
		}

		return c.JSON(fiber.Map{"Logout": "Success"})
	})

	if err := app.Listen("localhost:3001"); err != nil {
		panic(err)
	}
}

func createJWTToken(user data.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 60).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserId
	claims["FirstName"] = user.FirstName
	claims["PassportNumber"] = user.PassportNumber
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
