package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// --------
// infrastructure↓
// --------

var db *gorm.DB

// InitDB ...
func InitDB() *gorm.DB {
	dsn := "root:root@tcp(db:3306)/jwtsample?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// --------
// router↓
// --------

// InitRouting ...
func InitRouting(e *echo.Echo) {
	// 認証が不要なルーティング
	e.POST("/public", Public)
	e.POST("/signup", Signup)

	// 認証が必要なルーティングのグループ化
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(Config))

	// // 認証が必要なルーティング
	api.POST("/private", Private)
}

// --------
// Config↓
// --------

var signingKey = []byte(os.Getenv("SIGNINGKEY"))

// Config ...
var Config = middleware.JWTConfig{
	SigningKey: signingKey,
}

// --------
// handler↓
// --------

type SignupResponse struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

// Signup ...
func Signup(c echo.Context) error {
	user := User{}

	if err := c.Bind(&user); err != nil {
		return err
	}

	if user.Name == "" || user.Password == "" {
		return c.JSON(http.StatusCreated, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		})
	}

	u, err := FindUser(&user)
	if err != nil {
		c.JSON(http.StatusCreated, err)
	}
	if u != nil && u.ID != 0 {
		return c.JSON(http.StatusCreated, &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "Already exists",
		})
	}
	err = CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusCreated, err)
	}

	// ユーザー用トークン生成
	token := makeToken(user)

	return c.JSON(http.StatusCreated, SignupResponse{UserID: user.ID, Token: token})
}

// makeToken トークン生成
func makeToken(user User) string {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = user.ID
	claims["name"] = user.Name
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// 電子署名
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	return tokenString
}

// Public ...
func Public(c echo.Context) error {
	product := Product{}

	if err := c.Bind(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]string{"This is Public Function": product.Name})
}

// Private ...
func Private(c echo.Context) error {
	product := Product{}

	if err := c.Bind(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]string{"This is Private Function": product.Name})
}

// --------
// model↓
// --------

// User ...
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Product ...
type Product struct {
	Name string `json:"name"`
}

// --------
// data access↓
// --------

// CreateUser ...
func CreateUser(u *User) error {
	err := db.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

// FindUser ...
func FindUser(u *User) (*User, error) {
	user := &User{}
	err := db.Where("name = ? AND password = ?", u.Name, u.Password).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func main() {
	db = InitDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	e := echo.New()
	InitRouting(e)

	e.Start(":9010")
}
