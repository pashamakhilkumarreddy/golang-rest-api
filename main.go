package main

import (
	"fmt"

	routes "github.com/pashamakhilkumarreddy/golang-rest-api/api"
	app "github.com/pashamakhilkumarreddy/golang-rest-api/cmd"
	utils "github.com/pashamakhilkumarreddy/golang-rest-api/utils"
)

func main() {
	fmt.Println("Welcome Go Gorilla Mux REST API")
	a := &app.App{}
	a.Initialize()
	routes.InitializeRoutes(a)

	a.Run(utils.GetEnv("PORT", "8081"))
}
