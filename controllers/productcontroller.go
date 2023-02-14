package controllers

import (
	"encoding/json"
	"chat-ecomm/database"
	"chat-ecomm/entities"
	"log"
	"net/http"
	"strconv"

	"fmt"

	"github.com/gorilla/mux"
)

// func CreateProduct(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var product entities.Product             //new product variable
// 	json.NewDecoder(r.Body).Decode(&product) //decodes body of incoming json request and maps it to newly created variable
// 	database.Instance.Create(&product)       //Using GORM we create a new product by passing in the parsed product, which creates a new record in the products table
// 	json.NewEncoder(w).Encode(product)       //Returns the newly created product data back to the client
// }

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	productName := r.FormValue("exampleName1")
	priceValue := r.FormValue("examplePrice1")
	productPrice, _ := strconv.ParseFloat(priceValue, 64) //string to float64
	quantityValue := r.FormValue("exampleQuantity1")
	productQuantity, _ := strconv.ParseUint(quantityValue, 10, 64) //string to uint64, 10 represents decimal number
	productDescription := r.FormValue("exampleDescription1")
	product := entities.Product{Name: productName, Price: productPrice, Quantity: productQuantity, Description: productDescription}
	result := database.Instance.Create(&product)
	if result == nil {
		log.Fatal("Error inserting " + productName + " into database.")
		fmt.Fprintf(w, "Failed to add! Try again.")
	} else {
		fmt.Fprintf(w, "Added!")
	}
}

func checkIfProductExists(productId string) bool {
	var product entities.Product
	database.Instance.First(&product, productId)
	if product.ID == 0 {
		return false
	}
	return true
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productId := mux.Vars(r)["id"] //gets product id from query string of the request
	if checkIfProductExists(productId) == false {
		json.NewEncoder(w).Encode("Cannot find product!")
		return
	}
	var product entities.Product
	database.Instance.First(&product, productId) //product table is queried with product id
	//and fills in all product details to the newly created product variable
	json.NewEncoder(w).Encode(product)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []entities.Product
	database.Instance.Find(&products) //maps all available products from database to the products list variable
	json.NewEncoder(w).Encode(products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicaion/json")
	productId := mux.Vars(r)["id"]
	if checkIfProductExists(productId) == false {
		json.NewEncoder(w).Encode("Could not update as product not found!")
		return
	}
	var product entities.Product
	database.Instance.First(&product, productId)
	json.NewDecoder(r.Body).Decode(&product)
	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productId := mux.Vars(r)["id"]
	if checkIfProductExists(productId) == false {
		json.NewEncoder(w).Encode("Could not delete as product not found!")
		return
	}
	var product entities.Product
	database.Instance.Delete(&product, productId)           //GORM deletes the product by ID
	json.NewEncoder(w).Encode("Product deleted from cart!") //sent back to client
}
