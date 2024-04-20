package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	app "github.com/pashamakhilkumarreddy/golang-rest-api/cmd"

	. "github.com/pashamakhilkumarreddy/golang-rest-api/database/models" //lint:ignore ST1001 Ignore
	. "github.com/pashamakhilkumarreddy/golang-rest-api/utils"           //lint:ignore ST1001 Ignore
)

type AppHandlers struct {
	App *app.App
}

func (a *AppHandlers) Home(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"messages": "Hola Mundo!", "time": time.Now().String(), "database": "connected"})
}

func (a *AppHandlers) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	fmt.Printf("Fetching product with id %v\n", id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
	}

	p := Product{ID: id}
	if err := p.GetProduct(a.App.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "Product not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, p)
}

func (a *AppHandlers) GetProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching all products")
	start, _ := strconv.Atoi(r.FormValue("start"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if limit > 10 || limit < 1 {
		limit = 10
	}
	if start < 0 {
		start = 0
	}

	p := Product{}

	products, err := p.GetProducts(a.App.DB, start, limit)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, products)
}

func (a *AppHandlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating a new product")
	var p Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.CreateProduct(a.App.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, p)
}

func (a *AppHandlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Invalid product ID")
		return
	}
	fmt.Printf("Updating a product with id %v\n", id)
	var p Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.ID = id

	if err := p.UpdateProduct(a.App.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, p)
}

func (a *AppHandlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	fmt.Printf("Deleting a product with id %v\n", id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	p := Product{ID: id}
	if err := p.DeleteProduct(a.App.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusNoContent, map[string]string{"result": "success"})
}
