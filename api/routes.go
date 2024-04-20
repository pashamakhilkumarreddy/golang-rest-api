package routes

import (
	"github.com/pashamakhilkumarreddy/golang-rest-api/api/handlers"
	app "github.com/pashamakhilkumarreddy/golang-rest-api/cmd"
)

func InitializeRoutes(a *app.App) {
	h := &handlers.AppHandlers{App: a}

	apiV1Router := a.Router.PathPrefix("/api/v1").Subrouter()

	apiV1Router.HandleFunc("/", h.Home).Methods("GET")
	apiV1Router.HandleFunc("/health", h.Home).Methods("GET")
	apiV1Router.HandleFunc("/products", h.GetProducts).Methods("GET")
	apiV1Router.HandleFunc("/product", h.CreateProduct).Methods("POST")
	apiV1Router.HandleFunc("/product/{id:[0-9]+}", h.GetProduct).Methods("GET")
	apiV1Router.HandleFunc("/product/{id:[0-9]+}", h.UpdateProduct).Methods("PUT")
	apiV1Router.HandleFunc("/product/{id:[0-9]+}", h.DeleteProduct).Methods("DELETE")
}
