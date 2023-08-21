package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleSignIn(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"user":  "Lương",
		"email": "hoangkimluong192@gmail.com",
	})
}

func HandleSignUp(c echo.Context) error {
	type User struct {
		Emai     string `json:"email"`
		FullName string `json:"name"`
		Age      int `json:"age"`
	}

	user := User{
		Emai:     "hoangkimluong192@gmail.com",
		FullName: "Hoàng Kim Lương",
		Age:      21,
	}
	return c.JSON(http.StatusOK, user)
}