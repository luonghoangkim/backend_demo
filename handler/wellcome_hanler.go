package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Wellcome(c echo.Context) error {
	return c.String(http.StatusOK, "Wellcome to MyBlog")
}