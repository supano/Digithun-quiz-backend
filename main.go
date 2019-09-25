package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	ID       int    `json:"id" form:"id" query:"id"`
	Name     string `json:"name" form:"name" query:"name"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

type ResponseJson struct {
	Message string `json:"message" form:"message" query:"message"`
	Data    []User
}

type ResponseJsonSingle struct {
	Message string
	Data    *User
}

type DatabaseCon struct {
	Database *sql.DB
}

var dbCon DatabaseCon

func main() {
	e := echo.New()

	var err error
	dbCon.Database, err = sql.Open("mysql", "root:q1w2e3r4@tcp(127.0.0.1:3306)/digithun_quiz")
	if err != nil {
		e.Logger.Fatal(err)
	}

	defer dbCon.Database.Close()

	e.GET("/", hello)
	e.POST("/login", login)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.POST("/users", addUser)
	e.PATCH("/users/:id", editUser)

	e.Logger.Fatal(e.Start(":9000"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}

func login(c echo.Context) error {
	var request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusUnauthorized, "Email or password incorrect")
	}

	if request.Email == "supano1995@gmail.com" && request.Password == "q1w2e3r4"{
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Supachai Suthikeeree"
		claims["email"] = "supano1995@gmail.com"
		claims["exp"] = time.Now().Add(time.Hour).Unix()

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Email or password incorrect")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	} 

	return c.JSON(http.StatusUnauthorized, "Email or password incorrect")
}

func getUser(c echo.Context) error {

	var count int
	rowCount, err := dbCon.Database.Query("select count(id) from users where id = ?", c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseJson{
			Message: "Not Found",
		})
	}

	defer rowCount.Close()

	for rowCount.Next() {
		rowCount.Scan(&count)
	}

	if count < 1 {
		return c.JSON(http.StatusNotFound, ResponseJson{
			Message: "Not Found",
		})
	}

	rows, err := dbCon.Database.Query("select id, name, email from users where id = ?", c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseJson{
			Message: "Not Found",
		})
	}
	defer rows.Close()

	var (
		id    int
		name  string
		email string
		user  User
	)

	for rows.Next() {
		rows.Scan(&id, &name, &email)
		user = User{
			ID:    id,
			Name:  name,
			Email: email,
		}
	}

	return c.JSON(http.StatusOK, ResponseJsonSingle{
		Message: "Success",
		Data:    &user,
	})

}

func getUsers(c echo.Context) error {
	rows, err := dbCon.Database.Query("select * from users")
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseJson{
			Message: "Not Found",
		})
	}
	defer rows.Close()

	var (
		id       int
		name     string
		email    string
		password string
		users    []User
	)

	for rows.Next() {
		rows.Scan(&id, &name, &email, &password)
		users = append(users, User{
			ID:    id,
			Name:  name,
			Email: email,
		})
	}

	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, ResponseJson{
			Message: "Not Found",
		})
	}

	return c.JSON(http.StatusOK, ResponseJson{
		Message: "Success",
		Data:    users,
	})
}

func addUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	// users := append(users, &u)

	stmt, err := dbCon.Database.Prepare("INSERT INTO users(email, name, password) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(u.Name, u.Email, u.Password)
	defer stmt.Close()

	id, _ := result.LastInsertId()
	u.ID = int(id)

	if err != nil {
		return c.JSON(http.StatusOK, ResponseJsonSingle{
			Message: "Error",
		})
	}

	return c.JSON(http.StatusOK, ResponseJsonSingle{
		Message: "Success",
		Data:    u,
	})
}

func editUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	var (
		err error
	)
	rows, err := dbCon.Database.Query("UPDATE users SET name = ? WHERE id = ?", u.Name, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusOK, ResponseJsonSingle{
			Message: "Error",
		})
	}

	defer rows.Close()

	result, err := dbCon.Database.Query("SELECT id, name, email FROM users WHERE id = ?", c.Param("id"))

	defer result.Close()

	var (
		id    int
		name  string
		email string
		user  User
	)

	for result.Next() {
		result.Scan(&id, &name, &email)
		user = User{
			ID:    id,
			Name:  name,
			Email: email,
		}
	}

	return c.JSON(http.StatusOK, ResponseJsonSingle{
		Message: "Success",
		Data:    &user,
	})

}
