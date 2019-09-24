package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type user struct {
	Name     string `json:"name" form:"name" query:"name"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.POST("/users", func(c echo.Context) error {

		u := new(user)
		if err := c.Bind(u); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, u)
	})

	e.Logger.Fatal(e.Start(":9000"))
}

func addUser(name string, email string, password string) {
	fmt.Println(name)
}
