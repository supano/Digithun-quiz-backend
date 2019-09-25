package main

import (
	"database/sql"
	"log"
	"net/http"

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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.GET("/users", getUsers)

	e.POST("/users", addUser)

	e.Logger.Fatal(e.Start(":9000"))
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

	if len(users) > 0 {
		return c.JSON(http.StatusOK, ResponseJson{
			Message: "Success",
			Data:    users,
		})
	} else {
		return c.JSON(http.StatusNotFound, ResponseJson{
			Message: "Not Found",
		})
	}
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

	type res struct {
		Message string
		Data    *User
	}

	_, err = stmt.Exec(u.Name, u.Email, u.Password)
	if err != nil {
		return c.JSON(http.StatusOK, struct {
			Message string
		}{
			Message: "Error",
		})
	}

	return c.JSON(http.StatusOK, &res{
		Message: "Success",
		Data:    u,
	})
}
