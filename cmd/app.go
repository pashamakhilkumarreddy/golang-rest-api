package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

	user := os.Getenv("DATABASE_USER")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	dbPort, _ := strconv.Atoi(port)
	connStr := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	connectionString := fmt.Sprintf(connStr, host, dbPort, user, password, dbname)

	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Panic(err)
	}

	a.Router = mux.NewRouter()
}

func (a *App) Run(port string) {
	log.Printf("Server is up and running on port :%v\n", port)
	if err := http.ListenAndServe(":"+port, a.Router); err != nil {
		log.Fatal("Error starting application \n", err)
	}
}

