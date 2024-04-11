package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	user, password, dbname := os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME")

	connStr := "user=%s password=%s dbname=%s sslmode=disable"
	connectionString := fmt.Sprintf(connStr, user, password, dbname)

	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Panic(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(port string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %v", err)
	}
	log.Printf("Server is up and running on %v:%v\n", hostname, port)
	if err := http.ListenAndServe(":"+port, a.Router); err != nil {
		log.Fatal("Error starting application \n", err)
	}
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"messages": "Hola Mundo!", "time": time.Now().String()})
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	fmt.Printf("Fetching product with id %v\n", id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching all products")
	start, _ := strconv.Atoi(r.FormValue("start"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if limit > 10 || limit < 1 {
		limit = 10
	}
	if start < 0 {
		start = 0
	}

	p := product{}

	products, err := p.getProducts(a.DB, start, limit)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating a new product")
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid product ID")
		return
	}
	fmt.Printf("Updating a product with id %v\n", id)
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	fmt.Printf("Deleting a product with id %v\n", id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) initializeRoutes() {
	apiV1Router := a.Router.PathPrefix("/api/v1").Subrouter()

	apiV1Router.HandleFunc("/", a.home).Methods("GET")
	apiV1Router.HandleFunc("/health", a.home).Methods("GET")
	apiV1Router.HandleFunc("/products", a.getProducts).Methods("GET")
	apiV1Router.HandleFunc("/product", a.createProduct).Methods("POST")
	apiV1Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	apiV1Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	apiV1Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}
