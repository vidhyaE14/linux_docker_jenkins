package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductName    string  `json:"prodname"`
	ProductCategory string  `json:"prodcategory"`
	ProductPrice    float64 `json:"prodprice"`
	ProductStock    int     `json:"prodstock"`
}

var DB *gorm.DB
var err error

// func getDSN() string {
// 	// Retrieve the database credentials from environment variables
// 	username := os.Getenv("DB_USERNAME")
// 	password := os.Getenv("DB_PASSWORD")
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("DB_PORT")
// 	database := os.Getenv("DB_NAME")

// 	if username == "" || password == "" || host == "" || port == "" || database == "" {
// 		log.Fatal("One or more environment variables are not set")
// 	}

// 	// Construct the DSN string
// 	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		username, password, host, port, database)
// }

func initializeRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/products", GetProducts).Methods("GET")
	router.HandleFunc("/product/{id}", GetProduct).Methods("GET")
	router.HandleFunc("/products", CreateProduct).Methods("POST")
	router.HandleFunc("/product/{id}", UpdateProduct).Methods("PUT")
	router.HandleFunc("/product/{id}", DeleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe("0.0.0.0:8081", router)) 
}

func initializeMigration() {
	db_username:= os.Getenv ("DB_USERNAME")
	if db_username == ""{
		panic ("DB_USERNAME env variable is not set")
	}
	db_password:= os.Getenv ("DB_PASSWORD")
	if db_password == ""{
		panic ("DB_PASSWORD env variable is not set")
	}
	rds_endpoint:= os.Getenv ("RDS_ENDPOINT")
	if rds_endpoint == ""{
		panic ("RDS_ENDPOINT env variable is not set")
	}
	rds_name:= os.Getenv ("RDS_NAME")
	if rds_name == ""{
		panic ("rds_NAME env variable is not set")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",db_username,db_password,rds_endpoint,rds_name)
	//dsn := "admin:ywOO73jpZd@tcp(mysql-rds-1.c5eooawm67do.us-east-1.rds.amazonaws.com:3306)/inventorysystem?charset=utf8mb4&parseTime=True&loc=Local"                                  //getDSN() // Use environment variables to get DSN
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	DB.AutoMigrate(&Product{})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prods []Product
	DB.Find(&prods)
	json.NewEncoder(w).Encode(prods)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var prod Product
	result := DB.First(&prod, params["id"])
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(prod)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prod Product
	json.NewDecoder(r.Body).Decode(&prod)
	DB.Create(&prod)
	json.NewEncoder(w).Encode(prod)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var prod Product
	result := DB.First(&prod, params["id"])
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewDecoder(r.Body).Decode(&prod)
	DB.Save(&prod)
	json.NewEncoder(w).Encode(prod)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var prod Product
	result := DB.First(&prod, params["id"])
	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	DB.Delete(&prod, params["id"])
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("The product is deleted successfully!!!")
}

func main() {
	initializeMigration()
	initializeRouter()
	fmt.Println("Application is running")
}