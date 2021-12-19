package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// --------
// DB
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
// middleware↓
// --------

func InitMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderAccessControlAllowHeaders, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
}

// --------
// router↓
// --------

// InitRouting ...
func InitRouting(e *echo.Echo) {
	// 認証が不要なルーティング
	e.POST("/public", Public)
	e.POST("/login", Login)

	// 認証が必要なルーティングのグループ化
	auth := e.Group("/auth")

	// JWTの設定
	auth.Use(middleware.JWTWithConfig(Config))

	// // 認証が必要なルーティング
	auth.GET("/private", Private)
}

// --------
// JWTConfig
// --------

var signingKey = []byte(os.Getenv("SIGNINGKEY"))

// Config ...
var Config = middleware.JWTConfig{
	SigningKey: signingKey,
	Claims:     &jwtCustomClaims{}, // デフォルトは「jwt.MapClaims{}」だが、カスタムのstructを指定することも可能
}

// --------
// model↓
// --------

// LoginInfo ...
type LoginInfo struct {
	UID      string `json:"uid"`
	Password string `json:"password"`
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// --------
// handler↓
// --------

// Login ...
func Login(c echo.Context) error {
	loginInfo := &LoginInfo{}

	if err := c.Bind(loginInfo); err != nil {
		return err
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		Name:  "Jon Snow",
		Admin: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	if err != nil {
		fmt.Printf("SIGNINGKEY:%v\n", os.Getenv("SIGNINGKEY"))
		return err
	}

	return c.JSON(http.StatusOK, t)
}

// Public ...
func Public(c echo.Context) error {
	return c.JSON(http.StatusCreated, "This is Public Function")
}

// Private ...
func Private(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)       // JWTConfigに設定されたデフォルトのContextKey「user」を指定してToken取得
	claims := user.Claims.(*jwtCustomClaims) // Tokenからclaims取得
	name := claims.Name                      // claimsの設定内容取得
	return c.JSON(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	db = InitDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	e := echo.New()

	InitMiddleware(e)

	InitRouting(e)

	e.Start(":9010")
}
