package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	ID       int    `json:"id form:id query:id"`
	Name     string `json:"name" form:"name" query:"name"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

var (
	users = []*User{}
)

func main() {
	e := echo.New()

	db, err := sql.Open("mysql",
		"root:q1w2e3r4@tcp(127.0.0.1:3306)/digithun_quiz")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.GET("/users", func(c echo.Context) error {
		rows, err := db.Query("select * from users")
		if err != nil {
			log.Fatal(err)
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

		return c.JSON(http.StatusOK, users)
	})

	e.POST("/users", func(c echo.Context) error {

		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}

		// users := append(users, &u)

		stmt, err := db.Prepare("INSERT INTO users(email, name, password) VALUES(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(u.Name, u.Email, u.Password)
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":9000"))
}
