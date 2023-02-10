package data

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	UserId         int64 `gorm:"primaryKey"`
	Email          string
	PhoneNumber    string
	gender         string
	FirstName      string
	lastName       string
	PassportNumber int
	PasswordHash   string             `json:"-"`
	tokens         Unauthorized_token `gorm:"references:user_id"`
}

type Unauthorized_token struct {
	user_id    int
	token      string
	expiration time.Time
}

func CreateDBEngine() (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "postgres", "authServer")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Unauthorized_token{})

	return db, nil
}
