package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	main "github.com/pashamakhilkumarreddy/golang-rest-api"
)

var a main.App
var mutex sync.Mutex

func TestMain(m *testing.M) {
	a.Initialize()

	ensureTableExists()

	defer clearTable()

	code := m.Run()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatalf("Error in checking if table exists due to %v", err)
	}
}

func clearTable() {
	fmt.Println("Clearing table")
	mutex.Lock()
	defer mutex.Unlock()
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	fmt.Println("Clearing table complete")
}

const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS products
	(
		id SERIAL,
		name TEXT NOT NULL,
		price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
		CONSTRAINT products_pkey PRIMARY KEY (id)
	)
`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mutex.Lock()
	a.Router.ServeHTTP(rr, req)
	mutex.Unlock()

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, Got %d\n", expected, actual)
	}
}

func addProducts(count int) {
	fmt.Println("Adding products")
	if count < 1 {
		count = 1
	}

	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		fmt.Printf("Adding product with id %v\n", i+1)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a.DB.Exec("INSERT INTO products(name, price) VALUES ($1, $2)",
				"Product "+strconv.Itoa(i), (float64(i)+1.0)*10)
		}(i)
	}

	wg.Wait()
}

func TestEmptyTable(t *testing.T) {
	fmt.Println(strings.Repeat("-", 30))

	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/products", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	fmt.Println(strings.Repeat("-", 30))

	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/product/111", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got %s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	fmt.Println(strings.Repeat("-", 30))

	clearTable()

	var jsonStr = []byte(`{"name": "test product", "price": 11.22}`)

	req, _ := http.NewRequest("POST", "/api/v1/product", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product', Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22', Got '%v'", m["price"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected product id to be '1', Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	fmt.Println(strings.Repeat("-", 30))

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/api/v1/product/1", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	fmt.Println(strings.Repeat("-", 30))

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/api/v1/product/1", nil)
	response := executeRequest(req)

	var originalProduct map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/api/v1/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}
	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}
	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	fmt.Println(strings.Repeat("-", 30))

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/api/v1/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/v1/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/api/v1/products/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
