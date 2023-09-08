package main

import (
	"backend_myblog/db"
	"backend_myblog/handler"
	"backend_myblog/helper"
	"backend_myblog/log"
	repoimpl "backend_myblog/repository/repo_impl"
	"backend_myblog/router"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func init() {
	fmt.Println(" >>>>>> ", os.Getenv("APP_NAME"))
	fmt.Println("new logic")
	//os.Setenv("APP_NAME", "github")
	log.InitLogger(false)

}
func main() {
	fmt.Printf("main func")
	sql := &db.Sql{
		Host:     "host.docker.internal", //"localhost"
		Port:     5432,
		UserName: "postgres",
		PassWord: "Hoangkimluong192@",
		DbName:   "golang",
	}

	sql.Connect()
	defer sql.Close()

	e := echo.New()

	structValidator := helper.NewStructValidator()
	structValidator.RegisterValidate()

	e.Validator = structValidator

	userHandler := handler.UserHandler{
		UserRepo: repoimpl.NewUserRepo(sql),
	}

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
	}

	api.SetupRouter()
	e.Logger.Fatal(e.Start(":3000"))
}
