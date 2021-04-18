package main

import (
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
	e.POST("/auth", GetTokenHandler)

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

// Public ...
func Public(c echo.Context) error {
	post := post{}

	if err := c.Bind(&post); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, post)
}

// Private ...
func Private(c echo.Context) error {
	post := post{}

	if err := c.Bind(&post); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, post)
}

// GetTokenHandler ...
func GetTokenHandler(c echo.Context) error {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = "54546557354"
	claims["name"] = "taro"
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// 電子署名
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	return c.JSON(http.StatusCreated, tokenString)
}

type post struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
	URL   string `json:"url"`
}

// ---------------------------------------------------

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
