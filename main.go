package main

import (
	"backend_myblog/db"
	"backend_myblog/handler"
	"backend_myblog/log"
	repoimpl "backend_myblog/repository/repo_impl"
	"backend_myblog/router"
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func init(){
	fmt.Println("init package main")
	os.Setenv("APP_NAME", "github")
	log.InitLogger(false)
	
}
func main() {
	fmt.Printf("main func")
	sql := &db.Sql{
		Host: "localhost",
		Port: 5432,
		UserName: "postgres",
		PassWord: "Hoangkimluong192@",
		DbName: "golang",
	}

	sql.Connect()
	defer sql.Close()

	var email string
	err := sql.Db.GetContext(context.Background(), &email , "SELECT email FROM users WHERE email=$1", "hoangkimluong192@gmail.com")
	if err != nil {
		log.Error(err.Error())
	}
	
	
    e := echo.New() 

	userHandler := handler.UserHandler{
		UserRepo: repoimpl.NewUserRepo(sql),
	}
	
	api := router.API{
		Echo: e,
		UserHandler: userHandler,
	}

	api.SetupRouter()
    e.Logger.Fatal(e.Start(":3000"))
}


 
