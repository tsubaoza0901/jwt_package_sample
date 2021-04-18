package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

type post struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
	URL   string `json:"url"`
}

var public = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := &post{
		Title: "VueCLIからVue.js入門①【VueCLIで出てくるファイルを概要図で理解】",
		Tag:   "Vue.js",
		URL:   "https://qiita.com/po3rin/items/3968f825f3c86f9c4e21",
	}
	json.NewEncoder(w).Encode(post)
})

var private = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := &post{
		Title: "VGolangとGoogle Cloud Vision APIで画像から文字認識するCLIを速攻でつくる",
		Tag:   "Go",
		URL:   "https://qiita.com/po3rin/items/bf439424e38757c1e69b",
	}
	json.NewEncoder(w).Encode(post)
})

// package auth
// GetTokenHandler get token
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

	// JWTを返却
	w.Write([]byte(tokenString))
})

// JwtMiddleware check token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNINGKEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func main() {
	r := mux.NewRouter()
	// localhost:9010/publicでpublicハンドラーを実行
	r.Handle("/public", public)
	r.Handle("/private", JwtMiddleware.Handler(private)) // auth.JwtMiddleware.Handler(private)
	r.Handle("/auth", GetTokenHandler)                   // auth.GetTokenHandler

	//サーバー起動
	if err := http.ListenAndServe(":9010", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// // --------
// // model↓
// // --------

// // User ...
// type User struct {
// 	ID       int    `json:"id" gorm:"id"`
// 	Name     string `json:"name" gorm:"name"`
// 	Password string `json:"password" gorm:"password"`
// }

// // Todo ...
// type Todo struct {
// 	UID     int    `json:"uid"`
// 	ID      int    `json:"id" gorm:"praimaly_key"`
// 	Content string `json:"content"`
// }

// // // JwtCustomClaims ...
// // type JwtCustomClaims struct {
// // 	UID  int    `json:"uid"`
// // 	Name string `json:"name"`
// // 	jwt.StandardClaims
// // }

// // --------
// // infrastructure↓
// // --------

// var db *gorm.DB

// // InitDB ...
// func InitDB() *gorm.DB {
// 	dsn := "root:root@tcp(db:3306)/jwtsample?parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return db
// }

// // --------
// // router↓
// // --------

// // InitRouting ...
// func InitRouting(e *echo.Echo) {
// 	// 認証が不要なルーティング
// 	e.POST("/signup", Signup)
// 	e.POST("/login", Login)
// 	e.GET("/accessible", Accessible)

// 	// 認証が必要なルーティング
// 	api := e.Group("/api")
// 	api.Use(middleware.JWTWithConfig(Config))
// 	api.POST("/todos", AddTodo)
// 	// api.GET("/restricted", u.Restricted)
// }

// // --------
// // Config↓
// // --------

// var signingKey = []byte("secret")

// // Config ...
// var Config = middleware.JWTConfig{
// 	// Claims:     &JwtCustomClaims{},
// 	SigningKey: signingKey,
// }

// // --------
// // handler↓
// // --------

// // Signup ...
// func Signup(c echo.Context) error {
// 	user := User{}

// 	if err := c.Bind(&user); err != nil {
// 		return err
// 	}

// 	if user.Name == "" || user.Password == "" {
// 		return &echo.HTTPError{
// 			Code:    http.StatusBadRequest,
// 			Message: "invalid name or password",
// 		}
// 	}

// 	u, err := FindUser(&user)
// 	if err != nil {
// 		return err
// 	}
// 	if u != nil && u.ID != 0 {
// 		return &echo.HTTPError{
// 			Code:    http.StatusConflict,
// 			Message: "Already exists",
// 		}
// 	}
// 	err = CreateUser(&user)
// 	if err != nil {
// 		return err
// 	}

// 	// トークン作成
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["name"] = user.Name
// 	claims["admin"] = true
// 	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

// 	t, err := token.SignedString([]byte("secret"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusCreated, map[string]string{"token": t})
// }

// // Login ...
// func Login(c echo.Context) error {
// 	user := User{}

// 	if err := c.Bind(&user); err != nil {
// 		return err
// 	}
// 	u, err := FindUser(&user)
// 	if err != nil {
// 		return err
// 	}
// 	if u == nil {
// 		return &echo.HTTPError{
// 			Code:    http.StatusUnauthorized,
// 			Message: "User Not Found",
// 		}
// 	}

// 	if u.Name != user.Name || u.Password != user.Password {
// 		return &echo.HTTPError{
// 			Code:    http.StatusUnauthorized,
// 			Message: "invalid name or password",
// 		}
// 	}

// 	// claims := &JwtCustomClaims{
// 	// 	u.ID,
// 	// 	u.Name,
// 	// 	jwt.StandardClaims{
// 	// 		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
// 	// 	},
// 	// }

// 	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	// t, err := token.SignedString(signingKey)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// return c.JSON(http.StatusOK, map[string]string{
// 	// 	"token": t,
// 	// })
// 	return c.JSON(http.StatusCreated, "Login Success")
// }

// // Accessible ...
// func Accessible(c echo.Context) error {
// 	return c.String(http.StatusOK, "Accessible")
// }

// // AddTodo ...
// func AddTodo(c echo.Context) error {
// 	todo := Todo{}
// 	if err := c.Bind(&todo); err != nil {
// 		return err
// 	}

// 	if todo.Content == "" {
// 		return &echo.HTTPError{
// 			Code:    http.StatusBadRequest,
// 			Message: "invalid to or message fields",
// 		}
// 	}

// 	user := &User{}
// 	uid := CheckRegisteredUserFromToken(c)
// 	user, err := FindUser(user)
// 	if err != nil {
// 		return err
// 	}
// 	if user == nil {
// 		return c.JSON(http.StatusCreated, echo.ErrNotFound)
// 	}

// 	user.ID = uid
// 	CreateTodo(&todo)

// 	return c.JSON(http.StatusCreated, todo)
// }

// // CheckRegisteredUserFromToken 登録済み（Token発行済み）ユーザーかのチェック
// func CheckRegisteredUserFromToken(c echo.Context) int {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(*JwtCustomClaims)
// 	uid := claims.UID
// 	return uid
// }

// // --------
// // data access↓
// // --------

// // CreateUser ...
// func CreateUser(u *User) error {
// 	err := db.Create(&u).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // FindUser ...
// func FindUser(u *User) (*User, error) {
// 	user := &User{}
// 	err := db.Where("name = ? AND password = ?", u.Name, u.Password).First(user).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return user, nil
// }

// // CreateTodo  ...
// func CreateTodo(todo *Todo) error {
// 	err := db.Create(todo).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // --------
// // main.go↓
// // --------

// func main() {
// 	db = InitDB()

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer sqlDB.Close()

// 	e := echo.New()

// 	InitRouting(e)

// 	e.Start(":9010")
// }
