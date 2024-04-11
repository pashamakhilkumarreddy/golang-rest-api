package main

import "fmt"

func main() {
	fmt.Println("Welcome Go Gorilla Mux REST API")
	a := App{}
	a.Initialize()

	a.Run(getEnv("PORT", "8081"))
}
